import os
import whisper
import httpx
from fastapi import FastAPI, File, UploadFile, HTTPException
from fastapi.responses import JSONResponse
from fastapi.concurrency import run_in_threadpool
import asyncio
app = FastAPI()
model = whisper.load_model("base")
DISTRESS_WORDS = ["help", "emergency", "danger", "save me", "urgent"]
@app.post("/audio")
async def process_audio(file: UploadFile = File(...)):
    try:
        temp_audio_path = "temp_audio.wav"
        with open(temp_audio_path, "wb") as buffer:
            buffer.write(await file.read())
        transcription = await run_in_threadpool(transcribe_audio, temp_audio_path)
        distress_detected = any(word.lower() in transcription.lower() for word in DISTRESS_WORDS)
        if distress_detected:
            await send_distress_signal()
        os.remove(temp_audio_path)

        return JSONResponse(content={"transcription": transcription, "distress_detected": distress_detected})

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

def transcribe_audio(audio_path: str) -> str:
    result = model.transcribe(audio_path)
    return result["text"]

async def send_distress_signal():
    async with httpx.AsyncClient() as client:
        try:
            response = await client.post("http://localhost:5000/distress", json={"message": "Distress detected"})
            response.raise_for_status()
        except httpx.HTTPError as e:
            print(f"Failed to send distress signal: {e}")

async def periodic_ping():
    while True:
        print("Pinging...")
        await asyncio.sleep(3)

@app.on_event("startup")
async def startup_event():
    asyncio.create_task(periodic_ping())

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=5670)
