import {useState, useEffect} from "react"
import axios from "axios"

const AllTransactions = () => {
    // defines the structure of the returned data for debugging
    const [data, setData] = useState([
        {
            item: {
                budget_name: "",
                item_name: ""
            },
            transactions: [
                {
                    transaction_id: "",
                    name: "",
                    type: "",
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

    // stores the data the the user wanted to view
    const [dateString, setDateString] = useState(() => {
        const date_now = new Date()
        return `${date_now.getFullYear()}-${String(date_now.getMonth() + 1).padStart(2, "0")}`
    })

    const token = localStorage.getItem("token")

    // returns all transactions from all items of date set by the user
    useEffect(() => {
        const getTransaction = async () => {
            // get request to retrieve the transactions
            await axios.get(`http://localhost:8080/main/transactions/${date.year}/${date.month}`, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            }).then(response => {
                // store it in data
                setData(response.data)
            }).catch(error => {
            // return error message if there is any
            console.log(error.response?.data || error.message)
            })
        }
        getTransaction()
    }, [date, token])

    const handleMonthYearChange = (e) => {
        const date_string = e.target.value.split('-')
        setDateString(e.target.value)
        setDate(previousDate => ({
            ...previousDate,
            month: parseInt(date_string[1]),
            year: parseInt(date_string[0])
        }))
    }
    
    return (
        <div className="p-20 flex">
            <div>
                <input type="month" value={dateString} onChange={(e) => {
                    handleMonthYearChange(e);
                }}/>
            {Array.isArray(data) && data[0] !== null ? data.map((t, index) => (
                <div>
                    <h1>{t.item.budget_name}</h1>
                    <table>
                    <tr>
                        <th className="px-5">Name</th>
                        <th className="px-5">Type</th>
                        <th className="px-5">Amount</th>
                        <th className="px-5">Date</th>
                    </tr>
                    {Array.isArray(t.transactions) && t.transactions[0] !== null && t.transactions.map((t, index) => (
                        <tr key={index}>
                            <td className="px-5">{t.name}</td>
                            <td className="px-5">{t.type}</td>
                            <td className="px-5">£{t.amount}</td>
                            <td className="px-5">{`${t.date.substring(0, 10)}`}</td>
                        </tr>
                ))}                            
                    </table>

                </div>
            )) : <div>There are no transactions for this budget item</div>
        }                
            </div>

        </div>
    )
}

export default AllTransactions;