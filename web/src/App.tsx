"use client"

import './App.css'
import MainPage from './components/main/page'
import RootLayout from './layout'



function App() {

  return (
    <>
          <RootLayout>
              
              <div className='w-full p-5'>
                  <MainPage />
              </div>
          
          </RootLayout>
    </>
  )
}


export default App

