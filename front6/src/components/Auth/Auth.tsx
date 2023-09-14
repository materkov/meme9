import React from "react";
import * as styles from "./Auth.module.css";
import {authLogin, authRegister} from "../../api/api";
import {useGlobals} from "../../store/globals";
import {useNavigation} from "../../store/navigation";
import {Link} from "../Link/Link";

export function Auth() {
    const globals = useGlobals();
    const nav = useNavigation();

    const [email, setEmail] = React.useState('');
    const [password, setPassword] = React.useState('');

    const [error, setError] = React.useState('');
    const isReg = location.search === "?registration";

    if (location.search === "?logout") {
        document.cookie = "authToken=; Path=/; Expires=" + new Date(0).toGMTString();
        globals.setAuth({token: "", userId: "", userName: ""});
        nav.go("/");
        return null;
    }

    const vkAuthURL = "https://oauth.vk.com/authorize?client_id=7260220&response_type=code&v=5.131&redirect_uri=" + encodeURIComponent(location.origin + "/vk-callback");

    const onAuth = () => {
        if (!email || !password) {
            return;
        }
        setError('');

        if (isReg) {
            onRegister();
        } else {
            onLogin();
        }
    };

    const onLogin = () => {
        authLogin({email, password})
            .then(resp => {
                globals.setAuth(resp);

                document.cookie = "authToken=" + resp.token + "; Path=/";
                nav.go("/");
            })
            .catch((err) => {
                // TODO think about error codes
                if (err === "InvalidCredentials") {
                    setError('Invalid email or password')
                } else {
                    setError('Something wrong, please try again later')
                }
            })
    };

    const onRegister = () => {
        authRegister({email, password})
            .then(resp => {
                globals.setAuth(resp);

                document.cookie = "authToken=" + resp.token + "; Path=/";
                nav.go("/");
            })
            .catch((err) => {
                if (err === "EmailAlreadyRegistered") {
                    setError('This email already registered')
                } else {
                    setError('Something wrong, please try again later')
                }
            })
    }

    return <div>
        {isReg ?
            <h1>Registration</h1> :
            <h1>Login</h1>
        }

        <div className={styles.formRow}>
            <div className={styles.formLeft}>Email:</div>
            <input type="text" value={email} onChange={e => setEmail(e.target.value)}/>
        </div>

        <div className={styles.formRow}>
            <div className={styles.formLeft}>Password:</div>
            <input type="password" value={password} onChange={e => setPassword(e.target.value)}/>
        </div>

        <div className={styles.formRow}>
            <button type="button" onClick={onAuth}>Login</button>
        </div>

        {!isReg &&
            <div className={styles.formRow}>
                <a href={vkAuthURL}>Login via VK</a>
            </div>
        }

        <div className={styles.formRow}>
            {!isReg && <Link href="/auth?registration">Register new account</Link>}
            {isReg && <Link href="/auth">Login to existing account</Link>}
        </div>

        {error &&
            <div className={styles.error}>{error}</div>
        }
    </div>
}
