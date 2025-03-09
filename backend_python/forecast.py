import pandas as pd
import numpy as np
from scipy.stats import boxcox
from statsmodels.tsa.arima.model import ARIMA
from scipy.special import inv_boxcox

def forecast(transactions, months, p, d, q):
    # turns the transactions into a series
    transactions = pd.Series(transactions)
    model = ARIMA(transactions, order=(12, 1, 12)).fit()

    forecast = model.forecast(steps=months)
    
    return list(forecast), None

def mean_value(combined_transactions):
    #approximating the original and forecasted transactions as a polynomial
    x = np.arange(len(combined_transactions))
    coefficients = np.polyfit(x, combined_transactions, deg=7)
    fitted_poly = np.poly1d(coefficients)

    #integrate the polynomial
    integral = np.polyint(fitted_poly)
    first = x[0]
    last = x[-1]

    #calculate the mean value of the transactions and set it as recommended budget
    mean_value = (integral(last) - integral(first)) / (last - first)

    return mean_value