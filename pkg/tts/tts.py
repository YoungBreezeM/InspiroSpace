from gtts import gTTS

tts = gTTS(text="hello", lang='zh-tw')
tts.save("hello.mp3")