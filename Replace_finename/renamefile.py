import os
import re

def rename_files(directory):
    if not os.path.isdir(directory):
        return 'Указанная директория не существует или недоступна.'

    files = os.listdir(directory)
    renamed_files = []
    errors = []

    for file in files:
        lower_file = file.lower()
        if lower_file.endswith(('.pdf', '.epub')):
            file_type = '.pdf' if lower_file.endswith('.pdf') else '.epub'
            new_name = file.replace('dokumen.pub_', '')

            # Одно выражение для всех замен
            new_name = re.sub(r'(\d+-){0,4}\d{1,12}(?=\.{})|[\/\'\\]| '.format(re.escape(file_type)),
                              lambda m: '_' if m.group(0) == ' ' else '', new_name)

            new_path = os.path.join(directory, new_name)
            if new_name != file and not os.path.exists(new_path):
                try:
                    os.rename(os.path.join(directory, file), new_path)
                    renamed_files.append(new_name)
                except OSError as e:
                    errors.append(f'Ошибка при переименовании файла {file}: {e}')

    return renamed_files, errors if errors else 'Все файлы успешно переименованы.'


def create_registry(directory, renamed_files):
    # Получаем список всех PDF и EPUB файлов
    all_files = [file for file in os.listdir(directory) if file.lower().endswith(('.pdf', '.epub'))]
    total_files = len(all_files)

    with open(os.path.join(directory, 'registry.txt'), 'w') as f:
        # Запись общего количества файлов в шапку файла
        f.write(f"Общее количество файлов: {total_files}\n\n")

        # Запись реестра файлов
        for file in all_files:
            file_size_mb = os.path.getsize(os.path.join(directory, file)) / (1024 * 1024)  # Размер в мегабайтах
            f.write(f"{file}: {file_size_mb:.2f} MB\n")

if __name__ == "__main__":
    # Указываем директорию для поиска файлов
    directory_path = sys.argv[1]

    # Изменение имен файлов
    renamed_files = rename_files(directory_path)
    print("Имена файлов изменены:", renamed_files)

    # Создание реестра с подсчетом общего количества файлов
    create_registry(directory_path, renamed_files)
    print("Реестр файлов создан.")