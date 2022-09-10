import React from "react";
import styles from "./Loader.module.css";

export interface Props {
    width?: number;
    className?: string;
}

export class Loader extends React.Component<Props> {
    render() {
        const style = {
            'height': this.props.width || 50,
            'width': this.props.width || 50,
        };

        return (
            <div className={styles.Ring + " " + this.props.className}>
                <div className={styles.Ring__inner + " " + styles.Ring__inner1} style={style}/>
                <div className={styles.Ring__inner + " " + styles.Ring__inner2} style={style}/>
                <div className={styles.Ring__inner + " " + styles.Ring__inner3} style={style}/>
                <div className={styles.Ring__inner} style={style}/>
            </div>
        );
    }
}
