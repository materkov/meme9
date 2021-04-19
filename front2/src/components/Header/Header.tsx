import React from "react";
import styles from "./Header.module.css";
import {Link} from "../Link/Link";
import {FeedGetHeaderRequest, FeedGetHeaderResponse, HeaderRenderer} from "../../api/api2";
import {api} from "../../Api";

interface State {
    data?: HeaderRenderer;
}

export class Header extends React.Component {
    state: State = {};

    componentDidMount() {
        this.refreshData();
        setInterval(this.refreshData, 10 * 1000);
    }

    refreshData = () => {
        api<FeedGetHeaderRequest, FeedGetHeaderResponse>("meme.Feed", "GetHeader", {}).then(r => {
            this.setState({data: r.renderer});
        })
    }

    render() {
        let data = this.state.data;
        if (!data) {
            data = {
                userName: "",
                userAvatar: "",
                mainUrl: "/",
            }
        }

        return (
            <div className={styles.Header}>
                <Link href={data.mainUrl} className={styles.Logo}>meme</Link>

                <div className={styles.RightContainer}>
                    <div className={styles.Name}>{data.userName}</div>
                    <img className={styles.Avatar} alt="" src={data.userAvatar}/>
                </div>
            </div>
        );
    }
}
