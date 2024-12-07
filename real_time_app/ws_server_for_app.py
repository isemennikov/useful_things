#!/usr/bin/env python
pip install fastapi
from fastapi import FastAPI, WebSocket

app = FastAPI()

@app.websocket("/w")
async  def websocket_endpoint(webscoket: WebSocket):
    await webscoket.accept()
    while True:
        data = await webscoket.receive_text()
        await webscoket.send_text(f"You said: {data}")
