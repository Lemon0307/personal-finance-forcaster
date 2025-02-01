import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Line, Bar } from "react-chartjs-2";
import axios from "axios";

const Forecast = () => {

    const [budget, setBudget] = useState("")
    const [item, setItem] = useState("")
    const [months, setMonths] = useState(1)

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
            console.log(response.data)
        } catch (error) {
            alert(error.response?.data || error.message);
        }
    }

    const selectedBudgetData = budgetData.find(bd => bd.budget.budget_name === budget);
    
    return (
        <div className="p-20 flex">
            <div className="grid place-content-center">
                <h1>Forecast</h1>
                <div className="flex items-center">
                    <div>
                        {/* Budget Selection */}
                        <select onChange={handleSelectBudget} value={budget}>
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
                        {selectedBudgetData && selectedBudgetData.budget_items.length > 0 ? (
                            <select onChange={handleSelectItem} value={item}>
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
                            onChange={(e) => {setMonths(e.target.value); console.log(months)}}
                            />
                            <h1>months</h1>
                        </div>
                        <button onClick={() => ForecastTransactions()}>Forecast</button>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Forecast;