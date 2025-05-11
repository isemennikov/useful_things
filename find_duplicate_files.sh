#!/bin/bash

# Проверка наличия аргумента с путем к каталогу
if [ $# -eq 0 ]; then
    echo "Использование: $0 /путь/к/каталогу"
    exit 1
fi

DIRECTORY="$1"

# Проверка существования каталога
if [ ! -d "$DIRECTORY" ]; then
    echo "Ошибка: Каталог '$DIRECTORY' не существует!"
    exit 1
fi

echo "Сканирование каталога '$DIRECTORY' на наличие дубликатов файлов..."

# Временный файл для хранения хешей
TEMP_FILE=$(mktemp)

# Функция для отображения строки прогресса
show_progress() {
    local current=$1
    local total=$2
    local percentage=$((current * 100 / total))
    local completed=$((percentage / 2))
    local remaining=$((50 - completed))
    
    printf "\r[%-${completed}s%-${remaining}s] %d%%  (%d из %d файлов)" "$(printf '%0.s#' $(seq 1 $completed))" "$(printf '%0.s-' $(seq 1 $remaining))" "$percentage" "$current" "$total"
}

# Подсчет общего количества файлов для индикатора прогресса
echo "Подсчет количества файлов..."
TOTAL_FILES=$(find "$DIRECTORY" -type f | wc -l)
echo "Всего файлов: $TOTAL_FILES"

# Находим все обычные файлы и вычисляем их SHA128 хеши
echo "Вычисление SHA128 хешей..."

CURRENT=0
find "$DIRECTORY" -type f -print0 | while IFS= read -r -d $'\0' file; do
    # Вычисляем SHA128 хеш файла
    hash=$(shasum -a 128 "$file" | awk '{print $1}')
    echo "$hash|$file" >> "$TEMP_FILE"
    
    # Увеличиваем счетчик обработанных файлов
    CURRENT=$((CURRENT + 1))
    
    # Показываем прогресс каждые 10 файлов или для последнего файла
    if [ $((CURRENT % 10)) -eq 0 ] || [ "$CURRENT" -eq "$TOTAL_FILES" ]; then
        show_progress $CURRENT $TOTAL_FILES
    fi
done

# Завершаем строку прогресса
echo -e "\n\nАнализ хешей..."
echo "Поиск дубликатов..."

# Выводим информацию о статусе выполнения
echo "Сортировка и группировка результатов..."

# Подготовка переменных для подсчета статистики
duplicate_groups=0
total_duplicates=0

# Сортируем файл по хешам и ищем дубликаты
echo "Результаты поиска дубликатов:"
echo "=============================="

sort "$TEMP_FILE" | awk -F'|' '{
    hash = $1;
    file = $2;
    
    files_by_hash[hash] = files_by_hash[hash] ? files_by_hash[hash] "\n  - " file : "  - " file;
    count_by_hash[hash]++;
}
END {
    found = 0;
    for (hash in files_by_hash) {
        if (count_by_hash[hash] > 1) {
            print "Найдены дубликаты с хешем SHA128: " hash;
            print "Количество файлов: " count_by_hash[hash];
            print "Файлы:" files_by_hash[hash] "\n";
            found = 1;
        }
    }
    if (!found) {
        print "Дубликаты файлов не найдены.";
    }
}' > >(
    while IFS= read -r line; do
        echo "$line"
        
        # Подсчитываем количество групп дубликатов и общее количество дубликатов
        if [[ "$line" =~ ^"Найдены дубликаты" ]]; then
            duplicate_groups=$((duplicate_groups + 1))
        elif [[ "$line" =~ ^"Количество файлов: " ]]; then
            # Извлекаем число файлов
            file_count=${line#"Количество файлов: "}
            total_duplicates=$((total_duplicates + file_count))
        fi
    done
)

# Выводим статистику только если были найдены дубликаты
if [ "$duplicate_groups" -gt 0 ]; then
    echo "=============================="
    echo "Статистика:"
    echo "Найдено групп дубликатов: $duplicate_groups"
    echo "Общее количество файлов-дубликатов: $total_duplicates"
fi

# Удаление временного файла
rm "$TEMP_FILE"

echo "Сканирование завершено."
