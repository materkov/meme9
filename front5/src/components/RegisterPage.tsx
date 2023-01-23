import React from "react";
import styles from "./RegisterPage.module.css";
import {Link} from "./Link";
import {emailRegister} from "../store/actions/auth";
import {setRoute} from "../store/actions/route";

export function RegisterPage() {
    const [email, setEmail] = React.useState("test@email.com");
    const [password, setPassword] = React.useState("");
    const [error, setError] = React.useState("");

    const onRegister = () => {
        if (email === "" || password === "") {
            return;
        }

        setError('');

        emailRegister({
            email: email,
            password: password,
        }).then(() => {
            setRoute("/");
        }).catch((err) => {
            if (err === 'email already registered') {
                setError('Этот емейл уже зарегистрирован')
            }
        })
    }

    return <>
        <h2>Регистрация</h2>

        <input type="text" placeholder="Емейл" value={email} onChange={e => setEmail(e.target.value)}/>
        <input type="password" placeholder="Пароль" value={password} onChange={e => setPassword(e.target.value)}/>
        <button className={styles.submitBtn} onClick={onRegister}>Зарегистрироваться</button>

        {error &&
            <>
                <br/>
                <div className={styles.container}>{error}</div>
            </>
        }

        <br/><br/>
        <Link href={"/login"}>Вход</Link>
    </>;
}
