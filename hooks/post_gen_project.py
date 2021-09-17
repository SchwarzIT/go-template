import os
import shutil
from typing import List

GRPC_FILES: List[str] = ["api/proto", "tools.go", "buf.gen.yaml", "buf.yaml"]
OPENAPI_FILES: List[str] = ["api/openapi.v1.yml"]


def delete_empty_folder(path: str) -> None:
    """
    delete all empty folders

    Args:
        path (str, optional): Path the delete empty folders
    """
    # Iterate over the directory tree and check if directory is empty.
    for (dirpath, dirnames, filenames) in os.walk(path):
        if len(dirnames) == 0 and len(filenames) == 0:
            shutil.rmtree(path=dirpath)


def remove_objects(objects: List[str]) -> None:
    """
    Delete file or folder

    Args:
        objects (List[str]): list of files & folders to delete
    """
    for path in objects:

        # check if path is a file
        if os.path.isfile(path=path):
            os.remove(path=path)

        # if not no file, delete as folder
        else:
            shutil.rmtree(path=path, ignore_errors=True)


def remove_grpc_files() -> None:
    """
    remove grpc files & folders
    """
    remove_objects(objects=GRPC_FILES)


def remove_openapi_files() -> None:
    """
    remove openapi files
    """
    remove_objects(objects=OPENAPI_FILES)


if __name__ == "__main__":
    # remove grpc support, if selected
    if "{{ cookiecutter.grpc_enabled }}" == "yes":
        remove_openapi_files()
    else:
        remove_grpc_files()

    # delete empty folders
    delete_empty_folder(path=".")
