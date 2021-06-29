import hug
import json
from hug.middleware import CORSMiddleware
from numpy import False_
import pandas as pd
from sklearn.preprocessing import StandardScaler
from random import randint

from sklearn.decomposition import PCA
import matplotlib.pyplot as plt

url = "./server/DataSet-Covid-IA.csv"

fields = ['Fever', 'Tiredness', 'Dry-Cough', 'Difficulty-in-Breathing', 'Sore-Throat',
          'None_Sympton', 'Pains', 'Nasal-Congestion', 'Runny-Nose', 'Diarrhea', 'None_Experiencing',
          'Age_0-9', 'Age_10-19', 'Age_20-24', 'Age_25-59', 'Age_60+', 'Gender_Female', 'Gender_Male',
          'Severity_Mild', 'Severity_Moderate', 'Severity_None', 'Severity_Severe', 'Contact_Dont-Know', 'Contact_No', 'Contact_Yes']

df = pd.read_csv(url, skipinitialspace=True, usecols=fields)

features = ['Fever', 'Tiredness', 'Dry-Cough', 'Difficulty-in-Breathing', 'Sore-Throat',
            'Pains', 'Nasal-Congestion', 'Runny-Nose', 'Diarrhea',  'Gender_Female',
            'Contact_Dont-Know',  'Contact_Yes']

x = df.loc[:, features].values

x = StandardScaler().fit_transform(x)

pca = PCA(n_components=2)

principalComponents = pca.fit_transform(x)

# principalDF = pd.DataFrame(data=principalComponents, columns=[
# 'component1', 'component2'])

# principalDF.to_csv('server/CovidConPAC2.csv')

# Visualizar Data Set

api = hug.API(__name__)
api.http.add_middleware(CORSMiddleware(api))


@hug.post('/newdata')
def newdata(body):
    data = [
        body['Fever'],
        body['Tiredness'],
        body['DryCough'],
        body['DifficultyinBreathing'],
        body['SoreThroat'],
        body['Pains'],
        body['NasalCongestion'],
        body['RunnyNose'],
        body['Diarrhea'],
        body['Gender_Female'],
        body['Contact_DontKnow'],
        body['Contact_Yes']
    ]
    resp = transformation(data)
    datita = {
        "x": resp[0][0],
        "y": resp[0][1]
    }
    return datita


def transformation(data):
    newdat = pca.transform([[data[0], data[1], data[2], data[3], data[4],
                             data[5], data[6], data[7], data[8], data[9],
                             data[10], data[11]]])
    return newdat


# def api_astar():
#     data = [
#         body['Fever'],
#         body['Tiredness'],
#         body['DryCough'],
#         body['DifficultyinBreathing'],
#         body['SoreThroat'],
#         body['Pains'],
#         body['NasalCongestion'],
#         body['RunnyNose'],
#         body['Diarrhea'],
#         body['Gender_Female'],
#         body['Contact_DontKnow'],
#         body['Contact_Yes']
#     ]

#     result = transformation(data)

#     return jsonify(
#         x=result[0][0],
#         y=result[0][1]
#     )
