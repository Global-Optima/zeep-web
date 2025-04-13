import io
import json
import numpy as np
import face_recognition
from flask import Flask, request, jsonify

app = Flask(__name__)


@app.route("/health", methods=["GET"])
def health():
    """Проверка, жив ли сервис."""
    return jsonify({"status": "ok"}), 200


@app.route("/extract", methods=["POST"])
def extract_embedding():
    """
    Извлекает face-embedding из загруженного фото.
    Принимает multipart/form-data с полем "image".
    Возвращает JSON:
      {
        "success": bool,
        "embedding": [float, float, ...] или []
        "error": "текст_ошибки"
      }
    """
    if "image" not in request.files:
        return jsonify({"success": False, "error": "No image file provided"}), 400

    file = request.files["image"]
    image_bytes = file.read()
    if not image_bytes:
        return jsonify({"success": False, "error": "Empty image"}), 400

    try:
        image = face_recognition.load_image_file(io.BytesIO(image_bytes))
        locations = face_recognition.face_locations(image)
        if not locations:
            return jsonify({"success": False, "error": "No face detected"}), 400

        encodings = face_recognition.face_encodings(image, locations)
        if not encodings:
            return jsonify({"success": False, "error": "Failed to extract encoding"}), 400

        # Возьмём первый найденный encoding
        embedding = encodings[0].tolist()

        return jsonify({
            "success": True,
            "embedding": embedding
        }), 200

    except Exception as e:
        return jsonify({"success": False, "error": f"Exception: {str(e)}"}), 500


@app.route("/compare", methods=["POST"])
def compare_embedding():
    """
    Сравнивает face-encoding из присланного фото с embedding, который передаётся отдельно.

    Ожидает:
      - multipart/form-data: "image"
      - multipart/form-data: "embedding" (JSON-массив или строка)

    Возвращает JSON:
      {
        "match": bool,
        "distance": float,
        "error": "текст_ошибки"
      }
    """
    # 1) Получаем файл изображения
    if "image" not in request.files:
        return jsonify({"match": False, "error": "No image file provided"}), 400

    file = request.files["image"]
    image_bytes = file.read()
    if not image_bytes:
        return jsonify({"match": False, "error": "Empty image"}), 400

    # 2) Получаем embedding из формы
    embedding_str = request.form.get("embedding")  # строка, возможно JSON
    if not embedding_str:
        return jsonify({"match": False, "error": "No embedding provided"}), 400

    try:
        # Попробуем распарсить как JSON
        embedding_data = json.loads(embedding_str)
        known_encoding = np.array(embedding_data, dtype=np.float32)
    except Exception as e:
        return jsonify({"match": False, "error": f"Invalid embedding format: {e}"}), 400

    try:
        # Извлечь encoding из нового фото
        image = face_recognition.load_image_file(io.BytesIO(image_bytes))
        locations = face_recognition.face_locations(image)
        if not locations:
            return jsonify({"match": False, "error": "No face detected"}), 400

        encodings = face_recognition.face_encodings(image, locations)
        if not encodings:
            return jsonify({"match": False, "error": "Failed to extract encoding"}), 400

        new_encoding = encodings[0]

        # 3) Сравнение. compare_faces возвращает [True/False],
        #    а face_distance даёт численное расстояние.
        match_results = face_recognition.compare_faces([known_encoding], new_encoding)
        face_distances = face_recognition.face_distance([known_encoding], new_encoding)

        match = bool(match_results[0])
        distance = float(face_distances[0])

        return jsonify({
            "match": match,
            "distance": distance
        }), 200

    except Exception as e:
        return jsonify({"match": False, "error": f"Exception: {str(e)}"}), 500


if __name__ == "__main__":
    # Запуск Flask сервера на 8000 порту
    app.run(host="0.0.0.0", port=8000, debug=True)
