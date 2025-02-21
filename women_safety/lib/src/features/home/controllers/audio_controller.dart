import 'dart:async';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:record/record.dart';
import 'package:path_provider/path_provider.dart';
import 'dart:io';
import 'package:http/http.dart' as http;

class AudioController extends GetxController {
  final recorder =
      AudioRecorder(); // Updated: Use AudioRecorder instead of Record
  final RxBool isRecording = false.obs;
  Timer? _uploadTimer;
  String? _currentRecordingPath;

  // Local development server endpoint
  static const String _apiEndpoint = 'http://localhost:3000/audio';

  @override
  void onClose() {
    _uploadTimer?.cancel();
    recorder.dispose();
    super.onClose();
  }

  Future<void> toggleRecording() async {
    try {
      if (isRecording.value) {
        await stopRecording();
      } else {
        await startRecording();
      }
    } catch (e) {
      debugPrint('Recording error: $e');
    }
  }

  Future<void> startRecording() async {
    try {
      // Check and request permissions
      if (await recorder.hasPermission() == false) {
        debugPrint('Microphone permission denied');
        return;
      }

      // Get temporary directory for saving the recording
      final tempDir = await getTemporaryDirectory();
      _currentRecordingPath =
          '${tempDir.path}/audio_${DateTime.now().millisecondsSinceEpoch}.m4a';

      // Configure and start recording
      await recorder.start(
        RecordConfig(
          encoder: AudioEncoder.aacLc,
          bitRate: 128000,
          sampleRate: 44100,
        ),
        path: _currentRecordingPath ?? '',
      );

      isRecording.value = true;
      _startPeriodicUpload();
    } catch (e) {
      debugPrint('Start recording error: $e');
    }
  }

  Future<void> stopRecording() async {
    try {
      _uploadTimer?.cancel();
      await recorder.stop();
      isRecording.value = false;

      if (_currentRecordingPath != null) {
        await _uploadAudioFile(_currentRecordingPath!);
      }
    } catch (e) {
      debugPrint('Stop recording error: $e');
    }
  }

  void _startPeriodicUpload() {
    _uploadTimer = Timer.periodic(const Duration(seconds: 10), (timer) async {
      if (_currentRecordingPath != null) {
        await _uploadAudioFile(_currentRecordingPath!);
      }
    });
  }

  Future<void> _uploadAudioFile(String filePath) async {
    try {
      final file = File(filePath);
      if (!await file.exists()) return;

      final request = http.MultipartRequest('POST', Uri.parse(_apiEndpoint));
      request.files.add(
        await http.MultipartFile.fromPath(
          'audio',
          filePath,
          filename: 'audio_${DateTime.now().millisecondsSinceEpoch}.m4a',
        ),
      );

      final response = await request.send();
      debugPrint('Upload status: ${response.statusCode}');
    } catch (e) {
      debugPrint('Upload error: $e');
    }
  }
}
