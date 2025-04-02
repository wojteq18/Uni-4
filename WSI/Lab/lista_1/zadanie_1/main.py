import tensorflow as tf
import numpy as np
from sklearn.metrics import accuracy_score, precision_score, recall_score

#loading MNIST dataset
mnist = tf.keras.datasets.mnist
(x_train, y_train), (x_test, y_test) = mnist.load_data() #training and testing sets

# Normalizing pixel values to the range [0,1]
x_train = x_train / 255.0
x_test = x_test / 255.0

#Building the neutral network
model = tf.keras.models.Sequential([
    tf.keras.layers.Flatten(input_shape=(28, 28)),   # flattening the matrix to a vector of 784
    tf.keras.layers.Dense(128, activation='swish'),   # hidden layer
    tf.keras.layers.Dense(10, activation='softmax')  # 10 neurons in the output layer (digits 0-9), softmax makes probabilities
])

# Compiling the model
model.compile(optimizer='adam', #Adam + Nesterov momentum
              loss='sparse_categorical_crossentropy',
              metrics=['accuracy'])

# Training the model
model.fit(x_train, y_train, epochs=8)

# Making predictions
y_pred_probs = model.predict(x_test) #Calculates for each example the probability for the classes
y_pred = np.argmax(y_pred_probs, axis=1) # Choosing the class with the highest probability

# Calculating metrics
accuracy = accuracy_score(y_test, y_pred)
precision = precision_score(y_test, y_pred, average='macro') #Macro -> average for all classes
recall = recall_score(y_test, y_pred, average='macro')

print(f"accuracy: {accuracy:.4f}")
print(f"precision: {precision:.4f}")
print(f"recall: {recall:.4f}")

model.save("mnist_model.keras")
print("Model saved as mnist_model.keras")