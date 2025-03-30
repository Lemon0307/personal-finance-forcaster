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
    const [months, setMonths] = useState(1)
    const [items, setItems] = useState([
        {
            "forecasted_earning": [
                {
                    "Amount": 0,
                    "Month": 2,
                    "Year": 2025
                }
            ],
            "forecasted_spending": [
                {
                    "Amount": 0,
                    "Month": 2,
                    "Year": 2025
                }
            ],
            "item_name": "Meals",
            "net_cash_flow": 0,
            "recommended_budget": 0,
            "total_earning": [
                {
                    "Amount": 0,
                    "Month": 11,
                    "Year": 2024
                }
            ],
            "total_spending": [
                {
                    "Amount": 0,
                    "Month": 9,
                    "Year": 2024
                }
            ]
        }
    ])
    const [recommendedBudget, setRecommendedBudget] = useState([])
    const [budgetData, setBudgetData] = useState()

    const token = localStorage.getItem("token")

    useEffect(() => {
        // gets all budgets for the user to select to forecast
        const GetBudgets = async () => {
            await axios.get(`http://localhost:8080/main/budgets`, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            })
            .then(response => {
                // stores response data in a state
                setBudgetData(response.data)
            })
            .catch(error => {
                alert(error.response?.data || error.message);
            })
        }
        GetBudgets()
    }, [token])

    const handleSelectBudget = (e) => {
        const {value} = e.target
        setBudget(value)
    }

    const ForecastTransactions = async () => {
        // send request to the server to forecast transactions
        await axios.get(`http://localhost:8080/main/forecast/${months}/${budget}`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        }).then(response => {
            if (Array.isArray(response.data)) { // check if returned data is an array
                // stores forecasted data and recommended budgets to use later
                setItems(response.data)
                setRecommendedBudget(response.data.map(rb => ({
                    item_name: rb.item_name,
                    recommended_budget: rb.recommended_budget
                })))
                console.log(response.data)
            } else {
                alert("data isn't an array")
                setItems([])
            }
        }).catch(error => {
            alert(error.response?.data || error.message);
        })
    }

    const handleApplyBudget = async (item, recommended_budget) => {
        // apply the recommended budget for that specific item
        await axios.put(`http://localhost:8080/main/budgets/update_item/${budget}/${item}`, {
            budget_cost: recommended_budget
        }, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
        .then(response => {
            alert(response.data.Message) // alert response message
        })
        .catch(error => {
            alert(error.response?.data || error.message);
        })
    }

    const getRandomColour = () => {
        return `hsl(${Math.floor(Math.random() * 360)}, 70%, 50%)`;
    };

    // set x-axis values
    const labels = [...new Set(
            // extracts total spending and forecasted spending 
            // and puts them into a single array
            items.flatMap(item =>
                [...item.total_spending, ...item.forecasted_spending].map(s => `${s.Date}`)
            )
        )].sort((a, b) => {
        // sort the dates in ascending order
        const dateA = new Date(a + '-01');
        const dateB = new Date(b + '-01');
        return dateA - dateB;
    });
    
    console.log(labels)
    // make datasets for each item
    const datasets = items.flatMap(item => {
        const colour = getRandomColour();
        return [
        { // gather data for past spending in an item
            label: `${item.item_name} - Past Spending`,
            data: labels.map(label => {
                // extracts amount from total_spending array to use as y values
                const entry = item.total_spending.find(s => `${s.Date}` === label);
                return entry ? entry.Amount : null;
            }),
            borderColor: colour,
            backgroundColor: "rgba(0,0,0,0)",
            tension: 0.3,
            borderDash: [], // solid line
        },
        { // gather data for forecasted spending in an item
            label: `${item.item_name} - Forecasted Spending`,
            data: labels.map(label => {
                // extracts amount from forecasted_spending array to use as y values            
                const entry = item.forecasted_spending.find(s => `${s.Date}` === label);
                return entry ? entry.Amount : null;
            }),
            borderColor: colour,
            backgroundColor: "rgba(0,0,0,0)",
            tension: 0.3,
            borderDash: [5, 5], // dashed line
        }];
        });

    const forecastData = {
        labels,
        datasets,
    };

    const options = {
        responsive: true,
        plugins: {
            legend: {
                position: "top",
            }
        },
        datasets: {
            spanGaps: true,
        }
    };
    
        
    return (
        <div className="p-20 flex">
            <div className="grid place-content-center">
                <h1>Forecast</h1>
                <div className="flex">
                    <div>
                        {/* Budget Selection */}
                        <select onChange={(e) => handleSelectBudget(e)} value={budget}>
                            <option value="" disabled>Select Budget...</option>
                            {budgetData?.length > 0 ? (
                                budgetData.map((bd, index) => (
                                    <option key={index} value={bd.budget_name}>
                                        {bd.budget_name}
                                    </option>
                                ))
                            ) : (
                                <option disabled>No budgets available</option>
                            )}
                        </select>

                        {/* Budget Item Selection */}
                        <h1>Duration of forecast</h1>
                        <div className="flex">
                            <input 
                            type="number"
                            defaultValue={0}
                            value={months}
                            onChange={(e) => {setMonths(e.target.value)}}
                            min={1}
                            />
                            <h1>months</h1>
                        </div>
                        {budget && <button onClick={() => ForecastTransactions()}>Forecast</button>}
                        
                    </div>
                    <div className="flex-grow w-full">
                        <div style={{ height: "calc(100vh - 200px)" }}>
                        {items?.length > 0 && (
                        <Line
                            data={forecastData}
                            options={options}
                            className="w-full"
                        />
                        )}
                        </div>
                    </div>
                    {recommendedBudget && 
                    <div>
                        <h1>Recommended Budgets</h1>
                        {recommendedBudget.map((rb, index) => (
                            <div key={index}>
                                {rb.item_name} : {rb.recommended_budget}
                                <button 
                                onClick={() => handleApplyBudget(rb.item_name, rb.recommended_budget)}
                                className="bg-gray-100 p-2 rounded-xl"
                                    >Apply Budget</button>
                            </div>
                        ))}
                    </div>
                    }
                </div>

            </div>
        </div>
    )
}

export default Forecast;