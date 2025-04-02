import tensorflow as tf
import numpy as np
import cv2
import os
from sklearn.metrics import accuracy_score, precision_score, recall_score

#Loading model
model = tf.keras.models.load_model("../zadanie_1/mnist_model.keras")
print("Model loaded")

#Loading images
path = "MyDigits/"

filenames = [f for f in os.listdir(path) if f.endswith('.png')]

y_true = [] #true labels
y_pred = [] #predicted labels

for fname in filenames:
    filepath = os.path.join(path, fname)

    img = cv2.imread(filepath, cv2.IMREAD_GRAYSCALE) #Read image as grayscale

    img = cv2.resize(img, (28, 28))

    img = 255 - img #Invert colors

    img = img / 255.0

    img = img.reshape(1, 28, 28) #Checking 1 image (batch size)

    pred = model.predict(img) #Calculates for each example the probability for the classes
    pred_class = np.argmax(pred) #Choosing the class with the highest probability

    true_label = int(fname[0]) #Extracting true label from filename

    y_true.append(true_label)
    y_pred.append(pred_class)

print("Plik | Prawdziwa | Przewidziana")
for fname, true, pred in zip(filenames, y_true, y_pred):
    print(f"{fname} | {true} | {pred}")




accuracy = accuracy_score(y_true, y_pred)
precision = precision_score(y_true, y_pred, average='macro')
recall = recall_score(y_true, y_pred, average='macro')

print(f"Accuracy: {accuracy:.4f}")
print(f"Precision: {precision:.4f}")
print(f"Recall: {recall:.4f}")
