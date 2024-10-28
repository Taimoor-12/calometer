import { useState } from "react";
import { Link } from "react-router-dom";
import { http_post } from "./lib/http";
import Spinner from "./components/Spinner";
import "./Signup.css";

interface ValidationErrors {
  fullName?: string;
  username?: string;
  password?: string;
  confirmPassword?: string;
}

function Signup() {
  const apiUrl = process.env.REACT_APP_API_URL

  const [formData, setFormData] = useState({
    fullName: "",
    username: "",
    password: "",
    confirmPassword: "",
  });

  const [errors, setErrors] = useState<ValidationErrors>({})
  const [loading, setLoading] = useState(false)
  const [apiCallErrMsg, setApiCallErrMsg] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData({ ...formData, [name]: value })

    const error = validateField(name, value)
    setErrors({ ...errors, [name]: error })
  }

  const validateField = (name: string, value: string): string | undefined => {
    switch(name) {
      case "fullName":
        if (!value.trim()) {
          return "Full name is required"
        } else if (!/^[A-Za-z\s]+$/.test(value)) {
          return "Full name must contain only letters"
        }
        break
      case "username":
        if (!value.trim()) {
          return "Username is required"
        }
        break
      case "password":
        if (!value.trim()) {
          return "Password is required"
        } else if (value.length < 6) {
          return "Password must be at least 6 characters"
        }
        break
      case "confirmPassword":
        if (!value.trim()) {
          return "Confirm password is required"
        } else if (value !== formData.password) {
          return "Passwords do not match"
        }
        break
      default:
        break
    }

    return undefined
  }

  const handleSubmit  = async (e: React.FormEvent) => {
    e.preventDefault()
    
    const validationErrors: ValidationErrors = {}

    if (!formData.fullName.trim()) {
      validationErrors.fullName = "Full name is required"
    } else if (!/^[A-Za-z\s]+$/.test(formData.fullName)) {
      validationErrors.fullName = "Full name must contain only letters"
    }

    if (!formData.username.trim()) {
      validationErrors.username = "Username is required"
    }

    if (!formData.password.trim()) {
      validationErrors.password = "Password is required"
    } else if (formData.password.length < 6) {
      validationErrors.password = "Password must be at least 6 characters"
    }

    if (!formData.confirmPassword.trim()) {
      validationErrors.confirmPassword = "Confirm password is required"
    } else if (formData.password !== formData.confirmPassword) {
      validationErrors.confirmPassword = "Passwords do not match"
    }

    setErrors(validationErrors)


    if(Object.keys(validationErrors).length === 0) {
      setLoading(true)
      setApiCallErrMsg('')
      console.log("Form submitted", formData)
      const body = {
        name: formData.fullName,
        username: formData.username,
        password: formData.password
      }

      try {
        const resp = await http_post(`${apiUrl}/api/users/signup`, body)
        setLoading(false)
        const respCode = Object.keys(resp.code)[0]
        if (respCode !== "200" && respCode === "409") {
          setErrors({ ...errors, ["username"]: resp.code[respCode] })
          // setApiCallErrMsg(resp.code[respCode])
        } else if (respCode !== "200") {
          setApiCallErrMsg(resp.Code[respCode])
        }
      } catch (e) {
        setLoading(false)
        console.log(e)
        setApiCallErrMsg("Something went wrong, please try again")
      }
      
    }
  };

  return (
    <div className="main">
      <h1>SIGNUP</h1>

      <div className="parentDiv">
        <div className="formDiv">
          <form onSubmit={handleSubmit}>
            <div className="inputDiv">
              <p>Full Name</p>
              <input 
                type="text"
                name="fullName"
                value={formData.fullName} 
                placeholder="Enter your fullname"
                onChange={handleChange} 
              />
              { errors.fullName && <p className="error">{errors.fullName}</p> }
            </div>

            <div className="inputDiv">
              <p>Username</p>
              <input 
                type="text"
                name="username"
                value={formData.username} 
                placeholder="Enter your username"
                onChange={handleChange} 
              />
              { errors.username && <p className="error">{errors.username}</p> }
            </div>

            <div className="inputDiv">
              <p>Password</p>
              <div className="passwordContainer">
                <input 
                  type={showPassword ? "text" : "password"}
                  name="password"
                  value={formData.password} 
                  placeholder="Enter your password"
                  onChange={handleChange} 
                />
                <div className="showPasswordDiv" onClick={() => setShowPassword(!showPassword)}>
                  {showPassword ? <img src="/assets/eye.svg" alt="eye" /> : <img src="/assets/eye-off.svg" alt="eye-off" />}
                </div>
              </div>
              { errors.password && <p className="error">{errors.password}</p> }
            </div>

            <div className="inputDiv">
              <p>Confirm Password</p>
              <div className="passwordContainer">
                <input 
                  type={showConfirmPassword ? "text" : "password"}
                  name="confirmPassword"
                  value={formData.confirmPassword} 
                  placeholder="Enter the password again"
                  onChange={handleChange} 
                />
                <div className="showPasswordDiv" onClick={() => setShowConfirmPassword(!showConfirmPassword)}>
                  {showConfirmPassword ? <img src="/assets/eye.svg" alt="eye" /> : <img src="/assets/eye-off.svg" alt="eye-off" />}
                </div>
              </div>
              { errors.confirmPassword && <p className="error">{errors.confirmPassword}</p> }
            </div>
            
            <div className="registerBtn">
              <button type="submit" disabled={loading}>
                {loading ? <Spinner /> : "Register"}
              </button>
            </div>
            { apiCallErrMsg && <p className="error apiErrorMsg">{apiCallErrMsg}</p> }
          </form>
        
        <p className="loginPrompt">
          Already have an account? <span><Link to="/login">Login here</Link></span>
        </p>
        </div>
        <div className="imgDiv">
          <img src="/assets/healthy_habit.svg" alt="healthy habit" />
        </div>
      </div>
    </div>
  );
}

export default Signup;
