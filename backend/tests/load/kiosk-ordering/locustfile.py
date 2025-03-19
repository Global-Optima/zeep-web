# locustfile.py
from locust import HttpUser, TaskSet, task, between
import json

class OrderDrinksBehavior(TaskSet):
    def on_start(self):
        # Set cookies so that every request will have them.
        # Make sure these tokens are valid / up-to-date before testing.
        self.client.cookies.set(
            "EMPLOYEE_ACCESS_TOKEN",
            "eyJhbGciOiJIUzI1NiIsInRva2VuVHlwZSI6ImFjY2VzcyIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ6ZWVwLXdlYiIsImV4cCI6MTc0MTc2MzIwOCwiaWF0IjoxNzQxNjc2ODA4LCJpZCI6NCwicm9sZSI6IkJBUklTVEEiLCJ3b3JrcGxhY2VJZCI6MSwid29ya3BsYWNlVHlwZSI6IlNUT1JFIn0._rZQImXrUEYPLqTt9lDmw4V5FAq06-2YLc0wB5GVrFM"
        )
        self.client.cookies.set(
            "EMPLOYEE_REFRESH_TOKEN",
            "eyJhbGciOiJIUzI1NiIsInRva2VuVHlwZSI6InJlZnJlc2giLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJ6ZWVwLXdlYiIsImV4cCI6MTc0MjI4MTYwOCwiaWF0IjoxNzQxNjc2ODA4LCJpZCI6NCwicm9sZSI6IkJBUklTVEEiLCJ3b3JrcGxhY2VJZCI6MSwid29ya3BsYWNlVHlwZSI6IlNUT1JFIn0.TzBqu6blqiOsk-dyJy8Tj7XtZwfYN4x4xP9gRAxEaT4"
        )

    @task
    def order_drink(self):
        # Here is the body you provided
        data = {
            "customerName": "Мудрый Путешественник 906",
            "subOrders": [
                {
                    "storeProductSizeId": 17,
                    "quantity": 1,
                    "storeAdditivesIds": []
                }
            ]
        }
        # Make the POST request
        with self.client.post("/api/v1/orders",
                              data=json.dumps(data),
                              headers={"Content-Type": "application/json"},
                              catch_response=True) as response:
            if response.status_code != 200:
                response.failure(f"Failed to create order: {response.status_code} - {response.text}")

class KioskUser(HttpUser):
    tasks = [OrderDrinksBehavior]
    wait_time = between(1, 3)  # seconds between tasks

    # This is the base host for your requests
    host = "https://zeep.sytes.net"
