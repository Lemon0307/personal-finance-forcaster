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
    // const [monthFrom, setMonthFrom] = useState(1)
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

    const [budgetData, setBudgetData] = useState()

    const token = localStorage.getItem("token")

    useEffect(() => {
        const GetBudgets = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/main/budgets`, {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                })
                setBudgetData(response.data)
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

    const ForecastTransactions = async () => {
        await axios.get(`http://localhost:8080/main/forecast/${months}/${budget}`, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            }).then(response => {
                setItems(response?.data)
            }).catch(error => {
                alert(error.response?.data);
            })
    }

    // const handleApplyBudget = async () => {
    //     try {
    //         await axios.put(`http://localhost:8080/main/budgets/update_item/${budget}/${item}`, {
    //             budget_cost: recommendedBudget
    //         }, {
    //             headers: {
    //                 Authorization: `Bearer ${token}`
    //             }
    //         })
    //         alert("Successfully applied budget to item")
    //     } catch (error) {
    //         alert(error.response?.data || error.message);
    //     }
    // }

    const getRandomColor = () => {
        return `hsl(${Math.floor(Math.random() * 360)}, 70%, 50%)`;
    };

    console.log(items)

    const labels = [
        ...new Set(
          items.flatMap(item =>
            [...item.total_spending, ...item.forecasted_spending].map(s => `${s.Month}/${s.Year}`)
          )
        ),
      ].sort((a, b) => {
        // Sort dates in ascending order (e.g., "2/2025" before "3/2025")
        const [monthA, yearA] = a.split("/").map(Number);
        const [monthB, yearB] = b.split("/").map(Number);
        return yearA === yearB ? monthA - monthB : yearA - yearB;
      });

    // Create datasets for each item
    const datasets = items.flatMap(item => {
        const color = getRandomColor();
        return [
          {
            label: `${item.item_name} - Past Spending`,
            data: labels.map(label => {
              const entry = item.total_spending.find(s => `${s.Month}/${s.Year}` === label);
              return entry ? entry.Amount : null;
            }),
            borderColor: color,
            backgroundColor: "rgba(0,0,0,0)",
            tension: 0.3,
            borderDash: [], // Solid line for past data
          },
          {
            label: `${item.item_name} - Forecasted Spending`,
            data: labels.map(label => {
              const entry = item.forecasted_spending.find(s => `${s.Month}/${s.Year}` === label);
              return entry ? entry.Amount : null;
            }),
            borderColor: color,
            backgroundColor: "rgba(0,0,0,0)",
            tension: 0.3,
            borderDash: [5, 5], // Dashed line for forecast
          },
        ];
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
        },
        title: {
          display: true,
          text: "Forecasted Spending by Item",
        },
      },
    };
    
        
    return (
        <div className="p-20 flex">
            <div className="grid place-content-center">
                <h1>Forecast</h1>
                <div className="grid">
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
                            className="w-full h-full"
                        />
                        )}
                        </div>
                    </div>
                </div>
                {/* {items.length > 0 && 
                <div className="flex items-center">
                    <h1 className="px-5">Recommended Budget: £{recommendedBudget.toFixed(2)}</h1>
                    <button onClick={() => handleApplyBudget()} className="px-5">Apply Budget</button>
                </div>} */}

            </div>
        </div>
    )
}

export default Forecast;