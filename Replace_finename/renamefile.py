import os
import sys
import re


def rename_files(directory, extension):
    files = os.listdir(directory)
    renamed_files = []

    for file in files:
        if file.endswith(extension):
            # Replace spaces with underscores
            new_name = file.replace(' ', '_')
            # Replace /, \, and ' characters with '-'
            new_name = re.sub(r'[\/\'\\]', '-', new_name)
            # Remove "dokumet.pub_"
            new_name = new_name.replace('dokumet.pub_', '')
            # Remove 20 digits before extension
            new_name = re.sub(r'\d{20}(?=\.{})'.format(extension), '', new_name)

            # Rename file
            if new_name != file:
                os.rename(os.path.join(directory, file), os.path.join(directory, new_name))
                renamed_files.append(new_name)

    return renamed_files


def create_registry(directory):
    files = os.listdir(directory)
    files.sort()
    with open(os.path.join(directory, 'registry.txt'), 'w') as f:
        for file in files:
            file_size = os.path.getsize(os.path.join(directory, file))
            f.write(f"{file}: {file_size} bytes\n")


if __name__ == "__main__":
    # Получаем путь к директории и расширение файла из аргументов командной строки
    directory_path = sys.argv[1]
    file_extension = sys.argv[2]

    # Изменение имен файлов
    renamed_files = rename_files(directory_path, file_extension)
    print("Имена файлов изменены:", renamed_files)

    # Создание реестра
    create_registry(directory_path)
    print("Реестр файлов создан.")
