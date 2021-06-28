import flask
from flask_cors import CORS
from flask import request, jsonify
from numpy import False_
import pandas as pd
from sklearn.preprocessing import StandardScaler
from random import randint

from sklearn.decomposition import PCA
import matplotlib.pyplot as plt

url = "server/DataSet-Covid-IA.csv"

fields = ['Fever', 'Tiredness', 'Dry-Cough', 'Difficulty-in-Breathing', 'Sore-Throat',
          'None_Sympton', 'Pains', 'Nasal-Congestion', 'Runny-Nose', 'Diarrhea', 'None_Experiencing',
          'Age_0-9', 'Age_10-19', 'Age_20-24', 'Age_25-59', 'Age_60+', 'Gender_Female', 'Gender_Male',
          'Severity_Mild', 'Severity_Moderate', 'Severity_None', 'Severity_Severe', 'Contact_Dont-Know', 'Contact_No', 'Contact_Yes']

df = pd.read_csv(url, skipinitialspace=True, usecols=fields)

features = ['Fever', 'Tiredness', 'Dry-Cough', 'Difficulty-in-Breathing', 'Sore-Throat',
            'None_Sympton', 'Pains', 'Nasal-Congestion', 'Runny-Nose', 'Diarrhea', 'None_Experiencing',
            'Age_0-9', 'Age_10-19', 'Age_20-24', 'Age_25-59', 'Age_60+', 'Gender_Female', 'Gender_Male',
            'Severity_Mild', 'Severity_Moderate', 'Severity_None', 'Severity_Severe', 'Contact_Dont-Know', 'Contact_No', 'Contact_Yes']

x = df.loc[:, features].values

x = StandardScaler().fit_transform(x)

pca = PCA(n_components=2)

principalComponents = pca.fit_transform(x)

principalDF = pd.DataFrame(data=principalComponents, columns=[
    'component1', 'component2'])

principalDF.to_csv('server/CovidConPAC2.csv')

# Visualizar Data Set


app = flask.Flask(__name__)
CORS(app)

app.config["DEBUG"] = True


def transformation(data):
    newdat = pca.transform([[data[0], data[1], data[2], data[3], data[4],
                             data[5], data[6], data[7], data[8], data[9],
                             data[10], data[11], data[12], data[13], data[14],
                             data[15], data[16], data[17], data[18], data[19],
                             data[20], data[21], data[22], data[23], data[24]]])
    print(newdat)
    return newdat


@app.route('/api/v1/newdata', methods=['GET'])
def api_astar():
    data = [
        request.json['Fever'],
        request.json['Tiredness'],
        request.json['Dry-Cough'],
        request.json['Difficulty-in-Breathing'],
        request.json['Sore-Throat'],
        request.json['None_Sympton'],
        request.json['Pains'],
        request.json['Nasal-Congestion'],
        request.json['Runny-Nose'],
        request.json['Diarrhea'],
        request.json['None_Experiencing'],
        request.json['Age_0-9'],
        request.json['Age_10-19'],
        request.json['Age_20-24'],
        request.json['Age_25-59'],
        request.json['Age_60+'],
        request.json['Gender_Female'],
        request.json['Gender_Male'],
        request.json['Severity_Mild'],
        request.json['Severity_Moderate'],
        request.json['Severity_None'],
        request.json['Severity_Severe'],
        request.json['Contact_Dont-Know'],
        request.json['Contact_No'],
        request.json['Contact_Yes']
    ]
    result = transformation(data)

    return jsonify(
        x=result[0][0],
        y=result[0][1]
    )


app.run()
