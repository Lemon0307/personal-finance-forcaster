import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

const Login = () => {
  // Step management
  const [step, setStep] = useState(1);

  // User information state
  const [details, setDetails] = useState({
    username: "",
    email: "",
    password: "",
  });

  // Redirect function
  let redirect = useNavigate()

  // Security questions state
  const [securityQuestions, setSecurityQuestions] = useState([
    { question: "", answer: "" },
  ]);

  // Available security questions
  const questions = [
    "What is your favorite color?",
    "What is your mother's maiden name?",
    "What was your first pet's name?",
    "What city were you born in?",
  ];

  const handleDetailsChange = (e) => {
    const { name, value } = e.target;
    setDetails({ ...details, [name]: value });
  };

  const handleSecurityQuestionChange = (index, field, value) => {
    const updatedQuestions = [...securityQuestions];
    updatedQuestions[index][field] = value;
    setSecurityQuestions(updatedQuestions);
  };

  const addSecurityQuestion = () => {
    setSecurityQuestions([...securityQuestions, { question: "", answer: "" }]);
  };

  console.log(details)

  const handleSubmit = async (e) => {
    try {
      const userData = {
        user: {...details},
        security_questions: securityQuestions,
      };
      console.log(userData)
    const response = await axios.post("http://localhost:8080/login", userData)
    alert(response.data.Message)
    localStorage.setItem("token", response.data.Token)
    localStorage.setItem("username", userData.user.username)
    localStorage.setItem("email", userData.user.email)
    e.preventDefault()
    redirect('/')
    } catch (error) {
      alert(error.response.data)
    }
  };

  return (
    <div className="p-40">
      {step === 1 && (
        <div className="grid place-items-center">
          <h2>Login</h2>
          <input
            type="username"
            name="username"
            placeholder="Username"
            value={details.username}
            onChange={handleDetailsChange}
            className="py-2"
          />
          <input
            type="email"
            name="email"
            placeholder="Email"
            value={details.email}
            onChange={handleDetailsChange}
            className="py-2"
          />
          <input
            type="password"
            name="password"
            placeholder="Password"
            value={details.password}
            onChange={handleDetailsChange}
            className="py-2"
            />
          
          <button onClick={() => setStep(2)}>Next</button>
        </div>
      )}

      {step === 2 && (
        <div className="grid place-items-center">
          <h2>Security Questions</h2>
          {securityQuestions.map((sq, index) => (
            <div key={index} className="p-2">
              <select
                value={sq.question}
                onChange={(e) =>
                  handleSecurityQuestionChange(index, "question", e.target.value)
                }
              >
                <option value="" disabled>
                  Select a question
                </option>
                {questions.map((question, i) => (
                  <option key={i} value={question}>
                    {question}
                  </option>
                ))}
              </select>
              <input
                type="text"
                placeholder="Answer"
                value={sq.answer}
                onChange={(e) =>
                  handleSecurityQuestionChange(index, "answer", e.target.value)
                }
              />
            </div>
          ))}
          <div className="grid">
            <button onClick={addSecurityQuestion}>Add Another Question</button>
            {/* <button onClick={() => setStep(1)}>Back</button> */}
            <button onClick={handleSubmit}>Submit</button>            
          </div>
        </div>
      )}
    </div>
  );
};

export default Login;
