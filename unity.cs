using UnityEngine;
using System.IO.Ports;
using System;

// This replaces the Uduino functionality with direct serial communication. 
// Just attach this script to the same GameObject used with the Uduino controller, and 
// update references to call the new methods instead (commands maintain same structure used before)
public class BionicHandController : MonoBehaviour 
{
    SerialPort serialPort;
    public string portName = "COM3"; // Make configurable in Unity Inspector
    public int baudRate = 115200;

    void Start()
    {
        serialPort = new SerialPort(portName, baudRate);
        try {
            serialPort.Open();
        }
        catch (Exception e) {
            Debug.LogError($"Failed to open serial port: {e.Message}");
        }
    }

    // Replace the existing Uduino command functions with these
    public void SendCommand1(int prox02, int med02, int dist02, int lat02, 
                           int prox03, int med03, int dist03, int lat03)
    {
        if (serialPort != null && serialPort.IsOpen) {
            string command = $"1,{prox02},{med02},{dist02},{lat02},{prox03},{med03},{dist03},{lat03}\n";
            serialPort.WriteLine(command);
        }
    }

    public void SendCommand2(int prox04, int med04, int dist04, int lat04,
                           int prox05, int med05, int dist05, int lat05)
    {
        if (serialPort != null && serialPort.IsOpen) {
            string command = $"2,{prox04},{med04},{dist04},{lat04},{prox05},{med05},{dist05},{lat05}\n";
            serialPort.WriteLine(command);
        }
    }

    public void SendCommand3(int prox01, int dist01, int lat01)
    {
        if (serialPort != null && serialPort.IsOpen) {
            string command = $"3,{prox01},{dist01},{lat01}\n";
            serialPort.WriteLine(command);
        }
    }

    void OnApplicationQuit()
    {
        if (serialPort != null && serialPort.IsOpen) {
            serialPort.Close();
        }
    }
}
