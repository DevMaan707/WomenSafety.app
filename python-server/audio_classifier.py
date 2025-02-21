import os
import numpy as np
import pandas as pd
import librosa
import tensorflow as tf
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report, confusion_matrix
from tqdm import tqdm

class AudioDistressDetector:
    def __init__(self, sample_rate=22050, duration=5):
        self.sample_rate = sample_rate
        self.duration = duration
        self.model = None

    def extract_features(self, file_path):
        """Extract audio features using librosa"""
        try:
            audio, sr = librosa.load(file_path, sr=self.sample_rate, duration=self.duration)
            target_length = self.sample_rate * self.duration
            if len(audio) < target_length:
                audio = np.pad(audio, (0, target_length - len(audio)))
            else:
                audio = audio[:target_length]
            mfcc = librosa.feature.mfcc(y=audio, sr=sr, n_mfcc=13)
            chroma = librosa.feature.chroma_stft(y=audio, sr=sr)
            mel = librosa.feature.melspectrogram(y=audio, sr=sr)
            features = np.concatenate([
                mfcc.mean(axis=1),
                chroma.mean(axis=1),
                mel.mean(axis=1)
            ])
            return features

        except Exception as e:
            print(f"Error processing {file_path}: {str(e)}")
            return None

    def load_dataset(self, base_path):
        features = []
        labels = []

        for split in ['train', 'test']:
            split_path = os.path.join(base_path, split)

            normal_path = os.path.join(split_path, 'Normal')
            for audio_file in tqdm(os.listdir(normal_path)):
                if audio_file.endswith('.wav'):
                    file_path = os.path.join(normal_path, audio_file)
                    feature = self.extract_features(file_path)
                    if feature is not None:
                        features.append(feature)
                        labels.append(0)

            for category in ['Women', 'Child']:
                category_path = os.path.join(split_path, category)
                for audio_file in tqdm(os.listdir(category_path)):
                    if audio_file.endswith('.wav'):
                        file_path = os.path.join(category_path, audio_file)
                        feature = self.extract_features(file_path)
                        if feature is not None:
                            features.append(feature)
                            labels.append(1)
        return np.array(features), np.array(labels)

    def build_model(self, input_shape):
        """Create neural network model"""
        model = tf.keras.Sequential([
            tf.keras.layers.Dense(256, activation='relu', input_shape=(input_shape,)),
            tf.keras.layers.Dropout(0.3),
            tf.keras.layers.Dense(128, activation='relu'),
            tf.keras.layers.Dropout(0.3),
            tf.keras.layers.Dense(64, activation='relu'),
            tf.keras.layers.Dense(1, activation='sigmoid')
        ])

        model.compile(
            optimizer='adam',
            loss='binary_crossentropy',
            metrics=['accuracy']
        )

        return model

    def train(self, base_path, epochs=50, batch_size=32):
        print("Loading and preparing dataset...")
        X, y = self.load_dataset(base_path)
        X_train, X_val, y_train, y_val = train_test_split(
            X, y, test_size=0.2, random_state=42
        )
        print("Building and training model...")
        self.model = self.build_model(X_train.shape[1])
        history = self.model.fit(
            X_train, y_train,
            validation_data=(X_val, y_val),
            epochs=epochs,
            batch_size=batch_size,
            callbacks=[
                tf.keras.callbacks.EarlyStopping(
                    patience=5,
                    restore_best_weights=True
                )
            ]
        )
        return history

    def evaluate(self, X_test, y_test):
        predictions = (self.model.predict(X_test) > 0.5).astype(int)
        print("\nClassification Report:")
        print(classification_report(y_test, predictions))
        print("\nConfusion Matrix:")
        print(confusion_matrix(y_test, predictions))

    def predict(self, audio_path):
        features = self.extract_features(audio_path)
        if features is None:
            return None
        prediction = self.model.predict(features.reshape(1, -1))[0][0]
        return {
            'probability': float(prediction),
            'is_distress': bool(prediction > 0.5)
        }
