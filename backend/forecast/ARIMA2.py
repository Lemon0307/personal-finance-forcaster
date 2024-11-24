import numpy as np
import matplotlib.pyplot as plt

# Step 1: Use the provided time series data (e.g., monthly bills)
data = np.array([40, 47, 46, 44, 43, 46, 45, 47, 49, 55, 23, 45, 56])

# Time points (for simplicity, just use indices)
n = len(data)
time = np.arange(n)

# Plot the provided data
plt.figure(figsize=(10, 6))
plt.plot(time, data, label='Original Data')
plt.title('Original Time Series Data (Monthly Bills)')
plt.legend()
plt.show()

# Step 2: Differencing the data to make it stationary (I(1) part)
differenced_data = np.diff(data)  # First order differencing

# Plot the differenced data
plt.figure(figsize=(10, 6))
plt.plot(differenced_data)
plt.title('Differenced Data')
plt.show()

# Step 3: Estimate AR(1) parameter (phi)
def estimate_ar1(data):
    # Estimate the AR(1) parameter using the formula: phi = Cov(y_t, y_(t-1)) / Var(y_(t-1))
    lag_1 = data[:-1]
    y_t = data[1:]
    phi = np.corrcoef(lag_1, y_t)[0, 1]  # Correlation between y_t and y_(t-1) is equivalent to the AR(1) parameter
    return phi

phi_1 = estimate_ar1(differenced_data)
print(f"Estimated AR(1) coefficient (phi_1): {phi_1}")

# Step 4: Estimate MA(1) parameter (theta)
def estimate_ma1(data, phi_1):
    # Calculate the residuals (errors) after fitting the AR(1) model
    predicted_ar = np.roll(data, 1) * phi_1
    residuals = data - predicted_ar  # Residuals are the error term (y_t - predicted_y_t)
    residuals = residuals[1:]  # Exclude the first value because there's no prediction for it
    
    # Estimate the MA(1) parameter using the residuals
    theta_1 = np.corrcoef(residuals[:-1], residuals[1:])[0, 1]  # Correlation between residuals at time t-1 and t
    return theta_1

theta_1 = estimate_ma1(differenced_data, phi_1)
print(f"Estimated MA(1) coefficient (theta_1): {theta_1}")

# Step 5: Forecasting the next 10 values
def forecast_arima(data, phi_1, theta_1, n_forecast):
    # Forecast using the ARIMA model: y_t = phi_1 * y_(t-1) + theta_1 * residual_(t-1)
    predicted_values = []
    residuals = np.zeros(len(data))  # Initialize residual array
    
    for t in range(1, len(data)):
        # AR(1) part
        predicted_ar = phi_1 * data[t-1]
        
        # MA(1) part
        predicted_ma = theta_1 * residuals[t-1]
        
        # Forecast for the next time step
        predicted_value = predicted_ar + predicted_ma
        predicted_values.append(predicted_value)
        
        # Update residuals
        residuals[t] = data[t] - predicted_value

    # Forecasting future points beyond the data
    forecasted_values = []
    last_value = data[-1]  # Last observed value
    last_residual = residuals[-1]  # Last residual error
    
    for _ in range(n_forecast):
        # AR(1) and MA(1) for the next forecasted value
        forecast_ar = phi_1 * last_value
        forecast_ma = theta_1 * last_residual
        forecast_value = forecast_ar + forecast_ma
        
        # Add forecasted value to the list
        forecasted_values.append(forecast_value)
        
        # Update last_value and last_residual for the next iteration
        last_value = forecast_value + last_value  # Increment last observed value by forecast
        last_residual = forecast_value - forecast_ar  # Calculate new residual
    
    return np.array(predicted_values), np.array(forecasted_values)

# Forecast the next 10 months
n_forecast = 10
predicted_values, forecasted_values = forecast_arima(differenced_data, phi_1, theta_1, n_forecast)

# Step 6: Plotting the results
plt.figure(figsize=(12, 6))

# Original data and forecasted values
plt.plot(time, data, label='Original Data', color='blue')
plt.plot(np.arange(1, len(predicted_values) + 1), predicted_values, label='Predicted Data (ARIMA)', color='red')

# Add forecasted future data
forecast_time = np.arange(n, n + n_forecast)
plt.plot(forecast_time, forecasted_values, label='Forecasted Data', color='green')

plt.title('Original Data with ARIMA Forecast and Future Predictions')
plt.legend()
plt.show()

# Step 7: Invert differencing to return to original scale
# First, we need to add the previous value back to get the level of the forecasted data
forecasted_original_scale = np.concatenate(([data[0]], predicted_values)) + data[:-1]
forecasted_future_original_scale = forecasted_values + data[-1]  # Start forecast from last data point

# Plotting the forecasted original scale
plt.figure(figsize=(12, 6))
plt.plot(time, data, label='Original Data', color='blue')
plt.plot(np.arange(1, len(forecasted_original_scale) + 1), forecasted_original_scale, label='Forecasted Data (Original Scale)', color='red')

# Add forecasted future data on the original scale
plt.plot(forecast_time, forecasted_future_original_scale, label='Future Forecasted Data', color='green')

plt.title('Forecast and Future Predictions on Original Scale')

plt.show()
