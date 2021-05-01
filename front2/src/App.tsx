import ReactDOM from 'react-dom';
import React from "react";
import styles from "./App.module.css";
import {Router} from "./components/Router/Router";

function App() {
    return (
        <div className={styles.Background}>
            <div className={styles.App}>
                <Router/>
            </div>
        </div>
    )
}

ReactDOM.render(<App/>, document.getElementById('root'));
