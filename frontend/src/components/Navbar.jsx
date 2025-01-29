import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
// import io from "socket.io";

const Navbar = () => {
    const redirect = useNavigate()

    const logoutHandler = (e) => {
        localStorage.removeItem('token')
        localStorage.removeItem('email')
        localStorage.removeItem('username')
        redirect('/login')
    }

    const [currentBalance, setCurrentBalance] = useState(null)

    const token = localStorage.getItem("token")

    useEffect(() => {
        const socket = new WebSocket(`ws://localhost:8080/get_current_balance?token=${token}`);

        socket.onmessage = (event) => {
            const res = JSON.parse(event.data)
            setCurrentBalance(res.current_balance)
        }

        socket.onerror = (error) => {
            console.log(error)
        }
        
        return () => {
            if (socket.readyState === WebSocket.OPEN) {
                socket.close();
            }
        }
    }, [token])

    

    return (
        <div className="flex justify-around">
            <h1>Personal Finance Forecaster</h1>
            <div className="flex">
                {/* <select name="set_currency" id="" default="Set Currency...">
                    <option value="">Pound Sterling (£)</option>
                    <option value="">American Dollars ($)</option>
                </select> */}
                <h1 className="px-10">Current balance: £{currentBalance}</h1>
                <button  onClick={(e) => {e.preventDefault(); redirect('/')}}>Home</button>
                <button className="px-10" onClick={() => {logoutHandler()}}>Logout</button>                
            </div>

        </div>
    )
}

export default Navbar;