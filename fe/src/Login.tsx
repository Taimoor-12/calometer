import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { http_post, isRespDataWithHttpInfo } from "./lib/http";
import { toast } from "react-toastify";
import Spinner from "./components/Spinner";
import "./Login.css";

interface ValidationErrors {
  username?: string;
  password?: string;
}

function Login() {
  const apiUrl = process.env.REACT_APP_API_URL;

  const [formData, setFormData] = useState({
    username: "",
    password: "",
  });

  const [errors, setErrors] = useState<ValidationErrors>({});
  const [loading, setLoading] = useState(false);
  const [apiCallErrMsg, setApiCallErrMsg] = useState("");
  const [showPassword, setShowPassword] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });

    const error = validateField(name, value);
    setErrors({ ...errors, [name]: error });
  };

  const validateField = (name: string, value: string): string | undefined => {
    switch (name) {
      case "username":
        if (!value.trim()) {
          return "Username is required";
        }
        break;
      case "password":
        if (!value.trim()) {
          return "Password is required";
        } else if (value.length < 6) {
          return "Password must be at least 6 characters";
        }
        break;
      default:
        break;
    }

    return undefined;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const validationErrors: ValidationErrors = {};

    if (!formData.username.trim()) {
      validationErrors.username = "Username is required";
    }

    if (!formData.password.trim()) {
      validationErrors.password = "Password is required";
    } else if (formData.password.length < 6) {
      validationErrors.password = "Password must be at least 6 characters";
    }

    setErrors(validationErrors);

    if (Object.keys(validationErrors).length === 0) {
      setLoading(true);
      setApiCallErrMsg("");
      console.log("Form submitted", formData);
      const body = {
        username: formData.username,
        password: formData.password,
      };

      try {
        const resp = await http_post(`${apiUrl}/api/users/login`, body);
        setLoading(false);
        if (isRespDataWithHttpInfo(resp)) {
          const respCodeStr = Object.keys(resp.code)[0];
          const respCode: number = +respCodeStr;
          if (respCode === 409) {
            setErrors({ ...errors, ["username"]: resp.code[respCode] });
          } else if (respCode === 200) {
            toast.success("Login successful");
          } else {
            toast.error(resp.code[respCode]);
          }
        }
      } catch (e) {
        setLoading(false);
        console.log(e);
        setApiCallErrMsg("Something went wrong, please try again");
      }
    }
  };

  useEffect(() => {
    const loginUser = async () => {
      const resp = await http_post(`${apiUrl}/api/users/login`, {});
      if (isRespDataWithHttpInfo(resp)) {
        const respCodeStr = Object.keys(resp.code)[0];
        const respCode = +respCodeStr; // Convert string to number
        if (respCode === 200) {
          toast.success("Login successful"); //remove this later
          // navigate to dashboard/home.
        }
      }
    };

    loginUser();
  }, []);

  return (
    <div className="main">
      <h1>SIGNIN</h1>

      <div className="parentDiv">
        <div className="formDiv">
          <form onSubmit={handleSubmit}>
            <div className="inputDiv">
              <p>Username</p>
              <input
                type="text"
                name="username"
                value={formData.username}
                placeholder="Enter your username"
                onChange={handleChange}
              />
              {errors.username && <p className="error">{errors.username}</p>}
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
                <div
                  className="showPasswordDiv"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? (
                    <img src="/assets/eye.svg" alt="eye" />
                  ) : (
                    <img src="/assets/eye-off.svg" alt="eye-off" />
                  )}
                </div>
              </div>
              {errors.password && <p className="error">{errors.password}</p>}
            </div>

            <div className="loginBtn">
              <button type="submit" disabled={loading}>
                {loading ? <Spinner /> : "Login"}
              </button>
            </div>
            {apiCallErrMsg && (
              <p className="error apiErrorMsg">{apiCallErrMsg}</p>
            )}
          </form>

          <p className="signupPrompt">
            Don't have an account?{" "}
            <span>
              <Link to="/signup">Sign up here</Link>
            </span>
          </p>
        </div>
        <div className="imgDiv">
          <img src="/assets/login_illustration.svg" alt="login illustration" />
        </div>
      </div>
    </div>
  );
}

export default Login;
