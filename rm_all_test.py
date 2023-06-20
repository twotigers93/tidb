import os

def delete_test_files(directory):
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith("_test.go"):
                file_path = os.path.join(root, file)
                os.remove(file_path)
                print(f"Deleted file: {file_path}")

# 在当前目录中调用函数
delete_test_files(".")
