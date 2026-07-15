

function Hero(){
    const handleStart = () =>{
        alert("Bắt đầu sử dụng !")
    }
    return(
         <section id="home">
      <h1>
        Automate Your Marketing
        <br />
        Grow Your Business
      </h1>

      <p>
        Quản lý chiến dịch, khách hàng và tự động hóa marketing
        trên một nền tảng duy nhất.
      </p>

      <button onClick={handleStart}>
        Get Started
      </button>
    </section>
    )
}

export default Hero