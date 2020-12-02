

import { Layout } from 'antd';
const { Header} = Layout;

export function Navigation () {   
    return (
        <Header style={{"display":"block ruby", "zIndex": "1000", "position":"fixed", "left":0, "right":0, "top":0}}>
          <div className="logo"><h1 style={{"color":"white"}}>Watch2Gether</h1></div>
        </Header>
    )
}