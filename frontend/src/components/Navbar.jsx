import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import io from "socket.io";

const Navbar = () => {
    const redirect = useNavigate()

    const logoutHandler = (e) => {
        localStorage.removeItem('token')
        localStorage.removeItem('email')
        localStorage.removeItem('username')
        redirect('/login')
    }

    const [currentBalance, setCurrentBalance] = useState(null)
    const [socket, setSocket] = useState(null)

    const token = localStorage.getItem("token")

    useEffect(() => {
        const socket_connection = io('http://localhost:3000', {
            query: { token: token }  // Pass JWT token as a query parameter
        });
        setSocket(socket_connection)

        return () => {
            socket_connection.disconnect()
        }
    }, [token])

    

    return (
        <div className="flex justify-around">
            <h1>Personal Finance Forecaster</h1>
            <div>
                {/* <select name="set_currency" id="" default="Set Currency...">
                    <option value="">Pound Sterling (£)</option>
                    <option value="">American Dollars ($)</option>
                </select> */}
                <h1>{currentBalance}</h1>
                <button className="px-10" onClick={(e) => {e.preventDefault(); redirect('/')}}>Home</button>
                <button className="px-10" onClick={(e) => {logoutHandler()}}>Logout</button>                
            </div>

        </div>
    )
}

export default Navbar;