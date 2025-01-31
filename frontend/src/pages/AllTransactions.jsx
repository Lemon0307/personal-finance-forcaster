import {useState, useEffect} from "react"
import axios from "axios"

const AllTransactions = () => {
    const [data, setData] = useState([
        {
            budget_item: {
                budget_name: "",
                item_name: ""
            },
            transactions: [
                {
                    transaction_id: "",
                    transaction_name: "",
                    transaction_type: "",
                    amount: 0.00,
                    date: ""
                }
            ]
        }
    ])

    const [date, setDate] = useState({
        month: new Date().getMonth() + 1,
        year: new Date().getFullYear()
    })

    const token = localStorage.getItem("token")
    
    useEffect(() => {
        const getTransaction = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/main/transactions/${date.year}/${date.month}`, {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                })
                setData(response.data)
            } catch (error) {
                console.log(error.response?.data || error.message)
            }
        }
        getTransaction()
    }, [date, token])
    
    return (
        <div className="p-20 flex">
            <div className="px-10">
                <h1>Sort by:</h1>
                {/* <select value={sort} onChange={handleSort}>
                    <option value="name">Name</option>
                    <option value="type">Type</option>
                    <option value="amount">Amount</option>
                    <option value="date">Date</option>
                </select> */}
            </div>
            <div className="grid place-content-center">
                {Array.isArray(data) && data[0] !== null ?
                Array.isArray(data) && data[0] !== null ? data.map((t, index) => (
                    <div>
                        <h1>{t.budget_item.budget_name}</h1>
                        <table>
                        <tr>
                            <th className="px-5">Name</th>
                            <th className="px-5">Type</th>
                            <th className="px-5">Amount</th>
                            <th className="px-5">Date</th>
                        </tr>
                        {t.transactions.map((t, index) => (
                            <tr key={index}>
                                <td className="px-5">{t.transaction_name}</td>
                                <td className="px-5">{t.transaction_type}</td>
                                <td className="px-5">£{t.amount}</td>
                                <td className="px-5">{`${t.date.substring(0, 10)}`}</td>                                         
                            </tr>
                    ))}                            
                        </table>

                    </div>
                )) : <div>There are no transactions for this budget item</div>
                : <div></div>
            }
            </div>
        </div>
    )
}

export default AllTransactions;