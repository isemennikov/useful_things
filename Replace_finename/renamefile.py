import os
import re
import sys

import os
import re
import sys

def rename_files(directory):
    if not os.path.isdir(directory):
        return 'Указанная директория не существует или недоступна.', []

    files = os.listdir(directory)
    renamed_files = []
    errors = []
    file_type_pattern = re.compile(r'\.pdf$|\.epub$', re.IGNORECASE)

    for file in files:
        if file_type_pattern.search(file):
            file_extension = file_type_pattern.search(file).group()
            new_name = file.replace('dokumen.pub_', '')

            # Регулярное выражение для замены нежелательных символов и последовательностей
            new_name = re.sub(r'(\d+-){{0,4}}\d{{1,12}}|[/\'\\]|\s', lambda m: '_' if m.group(0).isspace() else '', new_name)

            # Убедимся, что имя файла оканчивается на правильное расширение
            if not new_name.endswith(file_extension):
                new_name += file_extension

            new_path = os.path.join(directory, new_name)
            if new_name != file:
                if not os.path.exists(new_path):
                    try:
                        os.rename(os.path.join(directory, file), new_path)
                        renamed_files.append(new_name)
                    except OSError as e:
                        errors.append(f'Ошибка при переименовании файла {file}: {e}')
                # Если файл с таким именем уже существует, пропускаем его
            else:
                errors.append(f'Файл {file} уже имеет требуемое имя и не был изменен.')

    return renamed_files, errors

def create_registry(directory, renamed_files):
    # Получение списка всех PDF и EPUB файлов
    all_files = [file for file in os.listdir(directory) if file.lower().endswith(('.pdf', '.epub'))]
    total_files = len(all_files)

    # Путь к файлу реестра
    registry_path = os.path.join(directory, 'registry.txt')
    with open(registry_path, 'w') as f:
        # Запись общего количества файлов в шапку файла
        f.write(f"Общее количество файлов: {total_files}\n\n")

        # Запись реестра файлов с их размерами
        for file in all_files:
            file_path = os.path.join(directory, file)
            file_size_mb = os.path.getsize(file_path) / (1024 * 1024)  # Размер в мегабайтах
            f.write(f"{file}: {file_size_mb:.2f} MB\n")

        # Если есть ошибки, добавляем их в конец файла
        if errors:
            f.write("\nОшибки:\n")
            for error in errors:
                f.write(f"{error}\n")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python renamefile.py <directory>")
        sys.exit(1)

    directory_path = sys.argv[1]
    renamed_files, errors = rename_files(directory_path)

    if renamed_files:
        print("Имена файлов изменены:", renamed_files)
    if errors:
        print("При переименовании возникли ошибки:\n - " + "\n - ".join(errors))

    create_registry(directory_path, renamed_files)
    print("Реестр файлов создан.")