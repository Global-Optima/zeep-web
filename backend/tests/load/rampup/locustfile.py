# locustfile.py
from locust import HttpUser, TaskSet, task, between

class KioskUserBehavior(TaskSet):
    def on_start(self):
        # Set cookies so that every request will have them.
        self.client.cookies.set(
            "EMPLOYEE_ACCESS_TOKEN",
            "eyJhbGciOiJIUzI1NiIsInRva2VuVHlwZSI6ImFjY2VzcyIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ6ZWVwLXdlYiIsImV4cCI6MTc0MTc2MzIwOCwiaWF0IjoxNzQxNjc2ODA4LCJpZCI6NCwicm9sZSI6IkJBUklTVEEiLCJ3b3JrcGxhY2VJZCI6MSwid29ya3BsYWNlVHlwZSI6IlNUT1JFIn0._rZQImXrUEYPLqTt9lDmw4V5FAq06-2YLc0wB5GVrFM"
        )
        self.client.cookies.set(
            "EMPLOYEE_REFRESH_TOKEN",
            "eyJhbGciOiJIUzI1NiIsInRva2VuVHlwZSI6InJlZnJlc2giLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJ6ZWVwLXdlYiIsImV4cCI6MTc0MjI4MTYwOCwiaWF0IjoxNzQxNjc2ODA4LCJpZCI6NCwicm9sZSI6IkJBUklTVEEiLCJ3b3JrcGxhY2VJZCI6MSwid29ya3BsYWNlVHlwZSI6IlNUT1JFIn0.TzBqu6blqiOsk-dyJy8Tj7XtZwfYN4x4xP9gRAxEaT4"
        )

        # EMPLOYEE_ACCESS_TOKEN=eyJhbGciOiJIUzI1NiIsInRva2VuVHlwZSI6ImFjY2VzcyIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ6ZWVwLXdlYiIsImV4cCI6MTc0MTc2MzIwOCwiaWF0IjoxNzQxNjc2ODA4LCJpZCI6NCwicm9sZSI6IkJBUklTVEEiLCJ3b3JrcGxhY2VJZCI6MSwid29ya3BsYWNlVHlwZSI6IlNUT1JFIn0._rZQImXrUEYPLqTt9lDmw4V5FAq06-2YLc0wB5GVrFM    ; 
        # EMPLOYEE_REFRESH_TOKEN=eyJhbGciOiJIUzI1NiIsInRva2VuVHlwZSI6InJlZnJlc2giLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJ6ZWVwLXdlYiIsImV4cCI6MTc0MjI4MTYwOCwiaWF0IjoxNzQxNjc2ODA4LCJpZCI6NCwicm9sZSI6IkJBUklTVEEiLCJ3b3JrcGxhY2VJZCI6MSwid29ya3BsYWNlVHlwZSI6IlNUT1JFIn0.TzBqu6blqiOsk-dyJy8Tj7XtZwfYN4x4xP9gRAxEaT4

    @task
    def load_kiosk(self):
        # Perform a GET request against the kiosk endpoint
        response = self.client.get("/kiosk")
        if response.status_code != 200:
            print("Request failed, status:", response.status_code)

class KioskUser(HttpUser):
    tasks = [KioskUserBehavior]
    # Wait times (in seconds) between task execution to simulate real user behavior
    wait_time = between(1, 5)
    
    # Host is the base URL
    host = "https://zeep.sytes.net"
