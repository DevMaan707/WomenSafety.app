from audio_classifier import AudioDistressDetector
detector = AudioDistressDetector()
history = detector.train('path/to/your/data/folder')
result = detector.predict('path/to/new/audio.wav')
print(f"Distress probability: {result['probability']:.2f}")
print(f"Is distress: {result['is_distress']}")
