from locust import HttpUser, task, between
import os

class ImageConverterUser(HttpUser):
    wait_time = between(1, 3)

    @task(3)
    def convert_lossy(self):
        """Task to test lossy conversion."""
        file_path = "test.png"
        if not os.path.exists(file_path):
            self.environment.runner.quit()
            raise Exception("Test image file 'test.png' not found!")
        with open(file_path, "rb") as image_file:
            files = {"image": ("test.png", image_file, "image/jpeg")}
            self.client.post("/convert?mode=lossy", files=files)

    @task(1)
    def convert_lossless(self):
        """Task to test lossless conversion."""
        file_path = "test.png"
        if not os.path.exists(file_path):
            self.environment.runner.quit()
            raise Exception("Test image file 'test.png' not found!")
        with open(file_path, "rb") as image_file:
            files = {"image": ("test.png", image_file, "image/jpeg")}
            self.client.post("/convert?mode=lossless", files=files)
