import os
import re
import sys

def rename_files(directory):
    files = os.listdir(directory)
    renamed_files = []

    for file in files:
        if file.lower().endswith(('.pdf', '.epub')):
            file_type = '.pdf' if file.lower().endswith('.pdf') else '.epub'

            # Замена пробелов и нежелательных символов
            new_name = file.replace(' ', '_')
            new_name = re.sub(r'[\/\'\\]', '-', new_name)
            new_name = new_name.replace('dokumet.pub_', '')

            # Удаление последовательности от 8 до 12 символов из чисел и знака "-" перед расширением файла
            new_name = re.sub(r'[-\d]{8,12}(?=\{})'.format(re.escape(file_type)), '', new_name)

            # Переименование файла
            if new_name != file:
                os.rename(os.path.join(directory, file), os.path.join(directory, new_name))
                renamed_files.append(new_name)

    return renamed_files

def create_registry(directory, renamed_files):
    total_files = len([file for file in os.listdir(directory) if file.lower().endswith(('.pdf', '.epub'))])
    with open(os.path.join(directory, 'registry.txt'), 'w') as f:
        # Запись общего количества файлов
        f.write(f"Общее количество файлов: {total_files}\n\n")
        # Запись реестра файлов
        for file in renamed_files:
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