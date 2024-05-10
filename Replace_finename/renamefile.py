import os
import re
import sys

import os
import re

def rename_files(directory):
    files = os.listdir(directory)
    renamed_files = []

    for file in files:
        if file.lower().endswith(('.pdf', '.epub')):
            file_type = '.pdf' if file.lower().endswith('.pdf') else '.epub'
            new_name = file

            # Удаление 'dokumen.pub_' и последовательности цифр и дефисов перед расширением файла
            new_name = new_name.replace('dokumen.pub_', '')
            new_name = re.sub(r'[-\d]{{8,12}}(?=\.{})'.format(re.escape(file_type)), '', new_name)

            # Замена пробелов на подчеркивания и удаление нежелательных символов
            new_name = new_name.replace(' ', '_')
            new_name = re.sub(r'[\/\'\\]', '-', new_name)

            # Переименование файла, если оно необходимо
            if new_name != file:
                os.rename(os.path.join(directory, file), os.path.join(directory, new_name))
                renamed_files.append(new_name)

    return renamed_files


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