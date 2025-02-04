import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Line, Bar } from "react-chartjs-2";
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend
} from "chart.js";
import axios from "axios";

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend);

const Forecast = () => {

    const [budget, setBudget] = useState("")
    const [item, setItem] = useState("")
    // const [monthFrom, setMonthFrom] = useState(1)
    const [months, setMonths] = useState(1)
    const [forecast, setForecast] = useState([])
    const [pastTransactions, setPastTransactions] = useState([{
        Month: 0,
        Year: 0,
        TotalAmount: 0.00
    }])

    const [budgetData, setBudgetData] = useState([
        {
            budget: {
                budget_name: ""
            },
            budget_items: [
                {
                    item_name: "",
                    budget_cost: 0.00,
                    description: "",
                    priority: 0
                }
            ]
        }
    ])

    const token = localStorage.getItem("token")
    const selectedBudgetData = budgetData.find(bd => bd.budget.budget_name === budget);

    useEffect(() => {
        const GetBudgets = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/main/budgets`, {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                })
                setBudgetData(response.data)
                // console.log(budget)
            } catch (error) {
                alert(error.response?.data || error.message);
            }
        }
        GetBudgets()
    }, [token])

    const handleSelectBudget = (e) => {
        const {value} = e.target
        setBudget(value)
    }

    const handleSelectItem = (e) => {
        const {value} = e.target
        setItem(value)
    }

    const ForecastTransactions = async () => {
        try {
            console.log(`http://localhost:8080/main/forecast/${months}/${budget}/${item}`)
            const response = await axios.get(`http://localhost:8080/main/forecast/${months}/${budget}/${item}`, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            })
            setForecast(response.data.forecast)
            setPastTransactions(response.data.total_transactions)
            // console.log(response.data)
        } catch (error) {
            alert(error.response?.data || error.message);
        }
    }

    const chartData = {
        labels: [...pastTransactions, ...forecast].map(entry => `Month ${entry.Month}/${entry.Year}`),
        datasets: [
            {
                label: "Past Transactions",
                data: pastTransactions.map(entry => entry.TotalAmount),
                borderColor: "blue",
                backgroundColor: "rgba(0, 0, 255, 0.2)",
                tension: 0.4,
            },
            {
                label: "Forecasted Transactions",
                data: [...new Array(pastTransactions.length).fill(null), ...forecast.map(entry => entry.TotalAmount)],
                borderColor: "red",
                backgroundColor: "rgba(255, 0, 0, 0.2)",
                tension: 0.4,
                borderDash: [5, 5]
            }
        ]
    };

    const options = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                position: "top",
            },
        },
        scales: {
            y: {
                beginAtZero: true,
            },
        },
    };
    
        
    return (
        <div className="p-20 flex">
            <div className="grid place-content-center">
                <h1>Forecast</h1>
                <div className="flex items-center">
                    <div>
                        {/* Budget Selection */}
                        <select onChange={(e) => handleSelectBudget(e)} value={budget}>
                            <option value="" disabled>Select Budget...</option>
                            {budgetData.length > 0 ? (
                                budgetData.map((bd, index) => (
                                    <option key={index} value={bd.budget.budget_name}>
                                        {bd.budget.budget_name}
                                    </option>
                                ))
                            ) : (
                                <option disabled>No budgets available</option>
                            )}
                        </select>

                        {/* Budget Item Selection */}
                        {selectedBudgetData && selectedBudgetData?.budget_items?.length > 0 ? (
                            <select onChange={(e) => handleSelectItem(e)} value={item}>
                                <option value="" disabled>Select Item...</option>
                                {selectedBudgetData.budget_items.map((item, index) => (
                                    <option key={index} value={item.item_name}>
                                        {item.item_name}
                                    </option>
                                ))}
                            </select>
                        ) : (
                            budget && <div>No budget items available</div>
                        )}
                        <h1>Duration of forecast</h1>
                        <div className="flex">
                            <input 
                            type="number"
                            defaultValue={0}
                            value={months}
                            onChange={(e) => {setMonths(e.target.value)}}
                            />
                            <h1>months</h1>
                        </div>
                        <button onClick={() => ForecastTransactions()}>Forecast</button>
                    </div>

                    {pastTransactions?.length > 0 && (
                        <div className="w-full h-screen">
                            <Line data={chartData} options={options}/>    
                        </div>    
                    )}
                </div>
            </div>
        </div>
    )
}

export default Forecast;