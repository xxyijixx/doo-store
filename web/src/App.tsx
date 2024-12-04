"use client"

import './App.css'
import MainPage from './components/main/page'
import RootLayout from './layout'



function App() {

  return (
    <>
          <RootLayout>
              
              <div className='w-full lg:pl-6 lg:pt-10 lg:pr-6  p-6 '>
                  <MainPage />
              </div>
          
          </RootLayout>
    </>
  )
}


export default App

