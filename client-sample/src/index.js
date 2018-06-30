import React from "react";
import ReactDOM from "react-dom";
import uuidv4 from "uuid/v4";
import Cookies from "js-cookie";
import PropTypes from "prop-types";
import url from "url";

const APIHost = "https://mining-gateway.appspot.com/";
const CoinhiveSiteKey = "99VOrUXDtysvhJCzuUzPicrKV41Fh9fg";

class App extends React.Component {
    state = {
        uid: null,
        miner: null,
        hashData: {
            hps: 0,
            total: 0,
            accepted: 0,
        },
        status: "stopping",
        lastAccepted: 0,
    };

    componentDidMount() {
        // miner setup
        const uidInCookie = Cookies.get("uid");
        const uid = uidInCookie || uuidv4();
        const miner = new window.CoinHive.User(CoinhiveSiteKey, uid);
        miner.on("accepted", () => {
            console.log("accepted");
            // set 2000ms delay to wait a transaction on Coinhive
            this.setState({lastAccepted: Date.now() + 3000});
        });
        this.setState({uid, miner});
        Cookies.set("uid", uid);

        // Update stats once per second
        setInterval(() => {
            const hashesPerSecond = miner.getHashesPerSecond();
            const totalHashes = miner.getTotalHashes();
            const acceptedHashes = miner.getAcceptedHashes();
            this.setState({hashData: {hps: hashesPerSecond, total: totalHashes, accepted: acceptedHashes}});
        }, 1000);
    }

    start = () => {
        this.state.miner.start();
        this.setState({status: "running..."});
    };
    stop = () => {
        this.state.miner.stop();
        this.setState({status: "stopping"});
    };

    render() {
        const {uid, status, lastAccepted} = this.state;
        const {hps, total, accepted} = this.state.hashData;

        return (
            <div>
                <h2>Paid API (Coinhive ver.) example</h2>

                <h3>API call</h3>
                <APICaller userID={uid} start={this.start} stop={this.stop} lastAccepted={lastAccepted}/>

                <h3>Manual Mining</h3>
                <div>status: {status}</div>
                <div>
                    <button onClick={this.start}>Start</button>
                    <button onClick={this.stop}>Stop</button>
                </div>
                <div>{hps} Hash/sec, TotalMined: {total}, <strong>TotalAccepted: {accepted}</strong></div>
                <div>Your ID (in Cookie): {uid}</div>

                <h3>Notice</h3>
                <div>We use cookies only to store your uid.</div>
                <div>This service is subject to finish without notice.</div>
            </div>
        );
    }
}

class APICaller extends React.Component {
    static propTypes = {
        userID: PropTypes.string,
        start: PropTypes.func.isRequired,
        stop: PropTypes.func.isRequired,
        lastAccepted: PropTypes.number.isRequired,
    };
    static defaultProps = {
        lastAccepted: 0,
    };

    state = {
        canTry: true, // if false, you have to mine before you try calling the API
        started: 0,
        callInfo: "",
    };

    static getDerivedStateFromProps(props, state) {
        if (state.started >= props.lastAccepted) {
            return null;
        }
        return Object.assign({}, state, {canTry: true, started: props.lastAccepted});
    }

    // this may fail
    callAPI = (endpoint) => async () => {
        this.setState({callInfo: "call start"});
        // try again...
        for (let i = 0; i < 50; i++) {
            this.setState({callInfo: `trying... ${i}`});
            if (this.state.canTry) {
                const res = await getResult(endpoint, this.props.userID);
                if (res) {
                    this.setState({callInfo: `call success - ${res}`});
                    this.props.stop();
                    return;
                }

                this.setState({callInfo: "failed to call, start mining..."});
                this.props.start();
                this.setState({canTry: false, started: Date.now()});
            }

            await new Promise(r => setTimeout(() => r(), 1000));
        }

        this.props.stop();
        this.setState({callInfo: "timeout. failed to call API... ><"});
        this.setState({canTry: true});
    };

    render() {
        const {callInfo} = this.state;
        return (
            <div>
                <div>
                    <button onClick={this.callAPI("/")}>call lightweight api (10 hash)</button>
                </div>
                <div>
                    <button onClick={this.callAPI("/heavy")}>call heavyweight api (250 hash)</button>
                </div>
                <div>info: {callInfo}</div>
            </div>
        );
    }
}

const getResult = async (endpoint, user) => {
    const target = url.resolve(APIHost, endpoint);
    const response = await fetch(target, {
        headers: {
            "X-Mining-Authorization": user,
        }
    });

    if (response.status !== 200) {
        console.log(response.status, response.json());
        return null;
    }

    return JSON.stringify(await response.json());
};

ReactDOM.render(<App/>, document.getElementById("index"));
