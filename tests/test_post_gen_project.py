import os
from pathlib import Path
from tempfile import TemporaryDirectory
from typing import List

from hooks.post_gen_project import delete_empty_folder, remove_objects


def test_delete_empty_folder() -> None:
    """
    test delete all empty folders
    """
    empty_folders: List[str] = ["empty_1", "empty_2/empty_3"]
    filled_folders: List[str] = ["filled_1"]
    files: List[str] = ["filled_1/file_1"]

    # create empty tmp folder
    with TemporaryDirectory() as temp_dir:
        # create folder
        folders: List[str] = empty_folders + filled_folders
        for folder in folders:
            os.makedirs(os.path.join(temp_dir, folder))
        # create files
        for file in files:
            Path(os.path.join(temp_dir, file)).touch()

        # delete empty folders
        delete_empty_folder(path=temp_dir)

        # check if all necessary folders was deleted
        for folder in empty_folders:
            assert os.path.exists(os.path.join(temp_dir, folder)) is False
        for folder in filled_folders:
            assert os.path.exists(os.path.join(temp_dir, folder)) is True


def test_remove_objects() -> None:
    """
    test delete file or folder
    """
    with TemporaryDirectory() as temp_dir:
        # create file
        file: str = os.path.join(temp_dir, "file")
        Path(file).touch()

        # create folder
        folder: str = os.path.join(temp_dir, "folder")
        os.makedirs(os.path.join(temp_dir, folder))

        # delete all files
        remove_objects(objects=[file, folder])

        # check if it was successfully deleted
        assert os.path.exists(file) is False
        assert os.path.exists(os.path.join(temp_dir, folder)) is False
