import './App.css';
import { useRef, useState, useEffect } from "react";
import { LuckyWheel } from '@lucky-canvas/react';
import Swal from "sweetalert2";
import withReactContent from "sweetalert2-react-content";

function App() {
  // const baseUrl = "http://127.0.0.1:81";
  const baseUrl = "";
  const LuckyRef = useRef();
  const blocks = [{ padding: '10px', background: '#869cfa' }];
  const wingif =
    "https://img.alicdn.com/imgextra/i3/917298378/O1CN01xehpw42BlAxlaOqJ3_!!917298378.gif";
  const bag =
    "https://img.alicdn.com/imgextra/i3/917298378/O1CN01UerpAW2BlAxjrJB70_!!917298378.png";
  const [prizes, setprizes] = useState([
    { fonts: [{ text: '加載中', top: '70%', fontSize: "12px" }], background: '#e9e8fe', imgs: [{ src: bag, top: "10%", width: "50%", height: "50%" }] },
    { fonts: [{ text: '加載中', top: '70%', fontSize: "12px" }], background: '#b8c5f2', imgs: [{ src: bag, top: "10%", width: "50%", height: "50%" }] },
    { fonts: [{ text: '加載中', top: '70%', fontSize: "12px" }], background: '#e9e8fe', imgs: [{ src: bag, top: "10%", width: "50%", height: "50%" }] },
    { fonts: [{ text: '加載中', top: '70%', fontSize: "12px" }], background: '#b8c5f2', imgs: [{ src: bag, top: "10%", width: "50%", height: "50%" }] },
    { fonts: [{ text: '加載中', top: '70%', fontSize: "12px" }], background: '#e9e8fe', imgs: [{ src: bag, top: "10%", width: "50%", height: "50%" }] },
    { fonts: [{ text: '加載中', top: '70%', fontSize: "12px" }], background: '#b8c5f2', imgs: [{ src: bag, top: "10%", width: "50%", height: "50%" }] },
    { fonts: [{ text: '加載中', top: '70%', fontSize: "12px" }], background: '#e9e8fe', imgs: [{ src: bag, top: "10%", width: "50%", height: "50%" }] },
    { fonts: [{ text: '加載中', top: '70%', fontSize: "12px" }], background: '#b8c5f2', imgs: [{ src: bag, top: "10%", width: "50%", height: "50%" }] },
  ]);
  const buttons = [
    { radius: '50px', background: '#617df2' },
    { radius: '45px', background: '#afc8ff' },
    {
      radius: '40px', background: '#ffb3eae0',
      pointer: true,
      fonts: [{ text: '開始\n抽獎', top: '-20px' }]
    },
  ];
  const mySwal = withReactContent(Swal);
  let prizesStatus = false;
  const luckData = useRef({});
  const swalpadding = "1.25rem";
  const swalbkgcolor = "rgb(254 240 255 / 50%)";
  const errorhandler = (error, icons) => {
    LuckyRef.current.stop();
    if (icons === "info") {
      mySwal.fire({
        padding: swalpadding,
        title: "Sorry",
        text: error,
        background: swalbkgcolor,
        confirmButtonColor: "#ea5455",
        imageUrl:
          "https://img.alicdn.com/imgextra/i3/917298378/O1CN01rYSA2P2BlAxmdQbE3_!!917298378.png",
        imageWidth: 150,
        footer: '<a href="line://">諮詢Anna老師</a>',
      });
      return;
    }
    mySwal.fire({
      padding: swalpadding,
      title: "提示",
      text: error,
      icon: icons,
      background: swalbkgcolor,
      confirmButtonColor: "#ea5455",
      footer: '<a href="line://">諮詢Anna老師</a>',
    });
  };
  const failPopup = () => {
    mySwal.fire({
      padding: swalpadding,
      title: luckData.current.Name,
      text: "未中獎",
      background: swalbkgcolor,
      confirmButtonColor: "#ea5455",
      imageUrl:
        "https://img.alicdn.com/imgextra/i2/917298378/O1CN01oYLUlV2BlAxl0nzaw_!!917298378.png",
      imageWidth: 150,
      footer: '<a href="line://">有問題？諮詢Anna老師</a>',
    });
  };
  const winPopup = () => {
    Swal.fire({
      padding: swalpadding,
      title: `恭喜中獎了\n<span style="color:red">>>${luckData.current.Name}<<</span>`,
      html: `<span>中獎LineID為: ${luckData.current.line}，趕快聯繫Anna老師領取禮品吧</span>`,
      imageUrl: luckData.current.Image,
      allowOutsideClick: false,
      confirmButtonColor: "#07b53b",
      showCancelButton: false,
      cancelButtonText: "關閉",
      imageWidth: 150,
      background: `${swalbkgcolor} url(${wingif})`,
      backdrop: `rgba(0,0,123,0.4) url("${wingif}") left top no-repeat`,
      footer: '<a href="line://">聯繫Anna老師領獎</a>',
    })
  }

  const luckedHandler = () => {
    if (!luckData.current.Win) {
      errorhandler("抽獎機會已用完，未中獎", "info");
      return;
    }
    mySwal.fire({
      padding: swalpadding,
      title: "Sorry,已經抽過了",
      html: `<span>中獎禮品為:  ${luckData.current.Name}<br/>
      中獎LineID為:  ${luckData.current.line}<br/>趕快聯繫Anna老師領取禮品吧</span>`,
      background: swalbkgcolor,
      confirmButtonColor: "#ea5455",
      imageUrl: luckData.current.Image,
      imageWidth: 150,
      footer: '<a href="line://">諮詢Anna老師</a>',
    });
  }
  const endHandler = () => {
    if (luckData.current.code === 400) {
      luckedHandler()
    } else if (luckData.current.code === 500) {
      errorhandler(luckData.current.message, "error")
    } else if (luckData.current.code === 200) {
      setTimeout(() => {
        LuckyRef.current.stop(luckData.current.Id);
      }, 2500);
    }
  }
  const setLuckyIndex = () => {
    // const p = fetch(baseUrl + "/luck/1111d9Xd/goodluck", {
    const p = fetch("./goodluck", {
      method: "POST",
      redirect: 'follow',
    });
    const t = new Promise((resolve, reject) => {
      setTimeout(() => {
        reject("請求超時！");
      }, 6000);
    });
    Promise.race([p, t])
      .then(response => response.json())
      .then(result => {
        if (result.code === 400 || result.code === 200) {
          luckData.current = { code: result.code, line: result.line, Id: result.data.Id, Image: result.data.Image, Name: result.data.Name, Win: result.data.Win };
        }
        if (result.code === 500) {
          luckData.current = { code: result.code, message: result.message }
        }
        if (result.code !== 200) {
          LuckyRef.current.stop();
        }
        endHandler();
      })
      .catch((error) => {
        console.error(error);
        LuckyRef.current.stop();
        errorhandler("你的網路好像出了點問題，請稍後再試", "error");
      });
  }
  const getPrizesList = async () => {
    if (prizesStatus) {
      return;
    }
    prizesStatus = true;
    // await fetch(baseUrl + "/luck/prizes", {
    await fetch(baseUrl + "/luck/prizes", {
      method: "POST",
      redirect: 'follow',
    })
      .then((res) => res.json())
      .then((result) => {
        if (result.code !== 200) {
          Swal.fire({
            padding: swalpadding,
            position: "center",
            icon: "error",
            title: "服務器似乎出現了點問題...",
            text: "暫時無法抽獎，請稍候再重試",
            background: swalbkgcolor,
            showConfirmButton: false,
            allowOutsideClick: false,
          });
          return;
        }
        const newprizes = prizes;
        result.data.forEach((item) => {
          if (newprizes[item.Id]) {
            newprizes[item.Id].fonts[0].text = item.Name;
            newprizes[item.Id].imgs[0].src = item.Image;
          }
        });
        setprizes(newprizes);
      })
      .catch((error) => {
        console.log(error)
        Swal.fire({
          padding: swalpadding,
          position: "center",
          icon: "error",
          title: "你的網路似乎有問題...",
          text: "暫時無法抽獎，請稍候再重試",
          background: swalbkgcolor,
          showConfirmButton: false,
          allowOutsideClick: false,
        });
      });
  };


  const start = function () {
    // 开始游戏
    LuckyRef.current.play();
    setLuckyIndex();
  };
  const end = function () { // 游戏停止时触发
    setTimeout(() => {
      if (luckData.current.code !== 200) return;
      if (luckData.current.Win) {
        winPopup()
      } else {
        failPopup()
      }
    }, 300)
  };
  useEffect(() => {
    if (!prizesStatus) {
      getPrizesList();
    }
  }, [prizesStatus]);
  return <LuckyWheel
    ref={LuckyRef}
    width="300px"
    height="300px"
    blocks={blocks}
    prizes={prizes}
    buttons={buttons}
    onStart={start}
    onEnd={end}
  ></LuckyWheel>

}

export default App;
