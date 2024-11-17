import numpy as np
import matplotlib.pyplot as plt

data = np.array([40, 47, 46, 44, 43, 46, 45, 47, 49, 55, 23, 45, 56])

n = len(data)
time = np.arange(n)

differenced_data = np.diff(data)

plt.figure(figsize=(10, 6))
plt.plot(differenced_data)
plt.title('Differenced Data')
plt.show()

def estimate_ar1(data):
    lag_1 = data[:-1]
    y_t = data[1:]
    phi = np.corrcoef(lag_1, y_t)[0, 1]
    return phi

phi_1 = estimate_ar1(differenced_data)
print(f"Estimated AR(1) coefficient (phi_1): {phi_1}")

def estimate_ma1(data, phi_1):
    predicted_ar = np.roll(data, 1) * phi_1
    residuals = data - predicted_ar
    residuals = residuals[1:]
    
    theta_1 = np.corrcoef(residuals[:-1], residuals[1:])[0, 1]
    return theta_1

theta_1 = estimate_ma1(differenced_data, phi_1)
print(f"Estimated MA(1) coefficient (theta_1): {theta_1}")

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

n_forecast = 10
predicted_values, forecasted_values = forecast_arima(differenced_data, phi_1, theta_1, n_forecast)

plt.figure(figsize=(12, 6))

plt.plot(time, data, label='Original Data', color='blue')
plt.plot(np.arange(1, len(predicted_values) + 1), predicted_values, label='Predicted Data (ARIMA)', color='red')

forecast_time = np.arange(n, n + n_forecast)
plt.plot(forecast_time, forecasted_values, label='Forecasted Data', color='green')

plt.title('Original Data with ARIMA Forecast and Future Predictions')
plt.legend()
plt.show()

forecasted_original_scale = np.concatenate(([data[0]], predicted_values)) + data[:-1]
forecasted_future_original_scale = forecasted_values + data[-1]

plt.figure(figsize=(12, 6))
plt.plot(time, data, label='Original Data', color='blue')
plt.plot(np.arange(1, len(forecasted_original_scale) + 1), forecasted_original_scale, label='Forecasted Data (Original Scale)', color='red')

plt.plot(forecast_time, forecasted_future_original_scale, label='Future Forecasted Data', color='green')

plt.title('Forecast and Future Predictions on Original Scale')
plt.legend()
plt.show()
