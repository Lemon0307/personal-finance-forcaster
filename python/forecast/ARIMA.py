import numpy as np
import matplotlib.pyplot as plt

# data = np.array([40, 47, 46, 44, 43, 46, 45, 47, 49, 55, 23, 45, 56])

def forecast_arima(data, phi_1, theta_1, n_forecast):
    predicted_values = []
    residuals = np.zeros(len(data))
    
    for t in range(1, len(data)):
        predicted_ar = phi_1 * data[t-1]
        predicted_ma = theta_1 * residuals[t-1]
        predicted_value = predicted_ar + predicted_ma
        predicted_values.append(predicted_value)
        residuals[t] = data[t] - predicted_value

    forecasted_values = []
    last_value = data[-1]
    last_residual = residuals[-1]
    
    for _ in range(n_forecast):
        forecast_ar = phi_1 * last_value
        forecast_ma = theta_1 * last_residual
        forecast_value = forecast_ar + forecast_ma
        
        forecasted_values.append(forecast_value)
        
        last_value = forecast_value + last_value
        last_residual = forecast_value - forecast_ar
    
    return np.array(predicted_values), np.array(forecasted_values)

def estimate_ar1(data):
    lag_1 = data[:-1]
    y_t = data[1:]
    phi = np.corrcoef(lag_1, y_t)[0, 1]
    return phi

def estimate_ma1(data, phi_1):
    predicted_ar = np.roll(data, 1) * phi_1
    residuals = data - predicted_ar
    residuals = residuals[1:]
    
    theta_1 = np.corrcoef(residuals[:-1], residuals[1:])[0, 1]
    return theta_1

def forecast_data(x):
    data = np.array(x)
    n = len(data)
    time = np.arange(n)

    differenced_data = np.diff(data)

    phi_1 = estimate_ar1(differenced_data)
    theta_1 = estimate_ma1(differenced_data, phi_1)
    n_forecast = 10
    predicted_values, forecasted_values = forecast_arima(differenced_data, phi_1, theta_1, n_forecast)
    forecasted_data = forecasted_values + data[-1]

    return forecasted_data

for element in forecast_data([432.45, 1124.10, 700.33, 1929.51, 1562.34, 810.23, 842.88, 454.76, 413.17, 356.22]):
    print(element)