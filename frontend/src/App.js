import { useState } from 'react';
import './App.css';

function App() {

  const [urls,setUrls] = useState("")
  const [pattern,setPattern] = useState("")
  const [isLoading,setIsLoading] = useState(false)
  const [result,setResult] = useState("")

  const scrapeUrls = async () => {
    setIsLoading(true)

    const urlsToScrape = urls.split("\n")
    const tags = pattern.split("->")

    const res = await fetch("http://localhost:5000/scrape", {
      method: "POST",
      headers: {
        "Content-type" : "json/application"
      },
      body: JSON.stringify({urls: urlsToScrape, tags}) 
    })

    const data = await res.json();

    setResult(data)
    
    setIsLoading(false)
  }

  return (
    <div className="w-full h-full flex justify-center items-center flex-col">
      <h1 className='text-5xl font-bold my-3'><span className='text-cyan-500'>Go</span>Scraper</h1>
      <div className='flex flex-col w-10/12'>
        <textarea className='text-xl my-1 border border-1 p-1' placeholder='https://en.wikipedia.org/wiki/Computer_programming' value={urls} onChange={(e) => setUrls(e.target.value)}>

        </textarea>
        <textarea className='text-xl my-1 border border-1 p-1' placeholder='div->h1->a' value={pattern} onChange={(e) => setPattern(e.target.value)}>

        </textarea>
        <button className='text-2xl font-bold px-2 py-1 bg-cyan-500 hover:bg-cyan-700 text-white transition-all' onClick={()=>scrapeUrls()}>SCRAPE</button>

      </div>
    </div>
  );
}

export default App;
