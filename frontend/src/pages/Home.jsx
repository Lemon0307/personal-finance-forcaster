import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

const Home = () => {

    const [transactions, setTransactions] = useState([])

    useEffect(() => {
        const getTransactions = async () => {
            let date = new Date()
        
            const token = localStorage.getItem('token')
            try {
                const response = await axios.get(`http://localhost:8080/main/transactions/${date.getFullYear()}/${date.getMonth() + 1}`,
                {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                }
            )
            setTransactions(response.data)
            } catch (error) {
                alert(error.response.data)
            }
        }
        getTransactions()
    }, [])

    const LineChart = () => {
        return (
            <div>
                <Line data={data} />
            </div>
        )
    }
    console.log(JSON.stringify(transactions))
    const username = localStorage.getItem('username')
    const [choice, setChoice] = useState()


    const choiceHandler = (e) => {
        setChoice(e.target.value)
    }

    return (
        <div>
            <h1>Welcome {username}</h1>
            <h1>Summary of transactions {choice}:</h1>
            <select
                value={choice}
                onChange={choiceHandler}
            >
                <option value="this week">This week</option>
                <option value="this month">This month</option>
            </select>
            <button>View budgets</button>
        </div>
    )
}

export default Home;