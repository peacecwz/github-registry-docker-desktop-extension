import * as React from "react";
import {Controls, Player} from "@lottiefiles/react-lottie-player";

const GithubLoader = () => {
    return <Player
        autoplay
        loop
        speed={2}
        src="https://assets9.lottiefiles.com/packages/lf20_S6vWEd.json"
        style={{height: '200px', width: '200px'}}
    >
        <Controls visible={false} buttons={['play', 'repeat', 'frame', 'debug']}/>
    </Player>
}

export default GithubLoader