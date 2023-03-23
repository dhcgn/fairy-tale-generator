# fairy-tale-generator

## Result/Ergebnis

Here is the result, with a really bad visual experience for fun :smirk:.

[![IMAGE ALT TEXT](https://img.youtube.com/vi/P6A515DSmDo/0.jpg)](https://www.youtube.com/watch?v=P6A515DSmDo "Amy und der Brückentroll - Demo mit Chat GTP und Amazon Polly")

## Run this tool

Get the Open AI API here https://openai.com/blog/openai-api and the Amazon AWS Key here https://aws.amazon.com.

```powershell
$env:OPENAI_API_KEY = "..."
$env:OPENAI_ORGANIZATION = "..."

$env:AWS_ACCESS_KEY_ID = "..."
$env:AWS_SECRET_ACCESS_KEY = "..."

.\fairy-tale-generator.exe
```

## English/Englisch :us:

A demo that showcases what Chat GPT can do for the story and Amazon Polly for the language.

Chat GPT creates the framework for a children's story based on a few inputs from the user and translates it into German. The model text-davinci-003 is used for this purpose. The result is a text.

Amazon Polly uses a neural language model to generate speech output from the text. This can be downloaded as an MP3.

In practice, an infinite number of stories can be generated this way. Limiting factors are the current limitations of Chat GPT in terms of complexity and length of the story, as well as the costs for using the AI SaaS, of course. Amazon Polly is currently considered to have the best language model, even though it still sounds artificial.

## German/Deutsch :de:

Dies ist nur eine Demo, die zeigt, wie Chat GPT eine Kindergeschichte und Amazon Polly die dazugehörige Sprache erzeugen kann.

Chat GPT erzeugt mit wenigen Benutzereingaben den Rahmen für eine Kindergeschichte und übersetzt diese ins Deutsche. Dazu wird das Modell text-davinci-003 verwendet. Das Ergebnis ist ein Text.

Amazon Polly verwendet ein neuronales Sprachmodell, um aus dem Text eine Sprachausgabe zu erzeugen. Diese wird als MP3-Datei heruntergeladen.

Praktisch kann eine unendliche Anzahl von Geschichten erzeugt werden, was durch die Anzahl der Parameter und die nicht-dererministische Natur des ML-Modells gerechtfertigt ist. Begrenzende Faktoren sind die aktuellen Grenzen von Chat GPT in Bezug auf Komplexität und Länge der Geschichten und natürlich die Kosten für die Nutzung der KI SaaS. Amazon Polly gilt derzeit als das beste Sprachmodell, auch wenn es noch künstlich klingt.

## Insights

### What is Amazon Polly?

Amazon Polly is a text-to-speech (TTS) service provided by Amazon Web Services (AWS). It is designed to convert written text into lifelike, natural-sounding speech. Amazon Polly uses advanced deep learning technologies to synthesize human-like voices for a wide range of applications, such as virtual assistants, conversational AI, e-learning platforms, audiobooks, news reading apps, and more.

### What is Chat GPT?

Chatbot GPT, or ChatGPT, is a conversational AI model based on the GPT (short for Generative Pre-trained Transformer) architecture, developed by OpenAI. The GPT series of models, including GPT-4, are large-scale language models that use deep learning techniques to understand and generate human-like text based on the input they receive.

### What is text-davinci-003?

"text-davinci-003" refers to a specific model offered by OpenAI as part of their API. It is based on the GPT-3 architecture, which is the third iteration of the Generative Pre-trained Transformer series. The "davinci" naming convention indicates that this model is designed to provide highly capable and sophisticated language understanding and generation.

The "003" in the name suggests it is a specific version or configuration within the "davinci" class of models. In general, GPT-3 models, including "text-davinci-003," are widely known for their ability to perform a variety of complex tasks, such as answering questions, writing content, summarization, translation, and even basic programming tasks, based on the input they receive.
