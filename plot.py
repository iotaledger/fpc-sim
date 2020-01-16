import numpy as np
import matplotlib.pyplot as plt

# - - - - most frequently changed output parameters - - -
# 0:a , 1:b , 2:beta , 3:k , 4:m , 5:l , 6:N , 7:p_0 , 8:q , 9:Adv_strategy
# 10:rateRandomness , 11:deltaWS , 12:gammaWS , 13:maxTermRound, 14:sZipf, 15:Type , 16:X , 17:Y
xcol = 3
xlabel = "k"
# xlim = [.0, 3.]
xlim = [1, 30]
# xscale = 'log'
xscale = 'linear'

# - - - - parameters unlikely to be changed - - -
ycol = 17  # last column in the csv file
folder = "data/"


def main():
    printDataRates()
    printDataFailures()
    printMeanTermRound()


def printDataRates():
    fig = plt.figure()
    filename = folder+'plot_ATI_rate'
    partPlot("Agreement", "AgreementRate",  ".", 10, filename, "blue")
    partPlot("Integrity", "IntegrityRate",   "+", 10, filename, "orange")
    partPlot("Termination", "TerminationRate",   "x", 10, filename, "green")
    plt.ylim([0, 1.05])
    plt.xscale(xscale)
    plt.xlim(xlim)
    plt.xlabel(xlabel)
    plt.ylabel("Rate")
    plt.legend(loc='best')
    plt.savefig(filename+'.eps', format='eps')
    plt.clf()


def printMeanTermRound():
    fig = plt.figure()
    filename = folder+'plot_EndRounds'
    partPlot("Last node", "MeanTerminationRound",
             ".", 10, filename, "magenta")
    partPlot("All nodes", "MeanLastRound",   "s", 5, filename, "darkgreen")
    plt.xlim(xlim)
    plt.xscale(xscale)
    plt.xlabel(xlabel)
    plt.ylim((0, 210))
    plt.ylabel("Mean Termination Round")
    plt.legend(loc='best')
    plt.savefig(filename+'.eps', format='eps')
    plt.clf()


def partPlot(type, file, marker, markersize, filename, color):
    x = loadDatafromRow(file, xcol)
    y = loadDatafromRow(file, ycol)
    x, y = sort2vecs(x, y)
    plt.plot(x, y, linestyle='dashed', color=color, linewidth=1)
    plt.plot(x, y, label=type, marker=marker,
             linestyle='none', color=color, markersize=markersize)
    np.savez(filename+"_"+type, x=x, y=y)


def printDataFailures():
    fig = plt.figure()
    filename = folder+'plot_ATI_failure'
    partPlot2("Agreement", "AgreementRate",  ".", 10, filename, "blue")
    partPlot2("Integrity", "IntegrityRate",   "+", 10, filename, "orange")
    partPlot2("Termination", "TerminationRate",   "x", 10, filename, "green")
    # plt.ylim([0, 1.05])
    plt.xscale(xscale)
    plt.yscale('log')
    plt.xlim(xlim)
    plt.xlabel(xlabel)
    plt.ylabel("Failure rate")
    plt.legend(loc='best')
    plt.savefig(filename+'.eps', format='eps')
    plt.clf()


def partPlot2(type, file, marker, markersize, filename, color):
    x = loadDatafromRow(file, xcol)
    y = loadDatafromRow(file, ycol)
    x, y = sort2vecs(x, y)
    plt.plot(x, 1-y, linestyle='dashed', color=color, linewidth=1)
    plt.plot(x, 1-y, label=type, marker=marker,
             linestyle='none', color=color, markersize=markersize)
    np.savez(filename+"_"+type, x=1-x, y=1-y)


def sort2vecs(x, y):
    i = np.argsort(x)
    x = x[i]
    y = y[i]
    return x, y


def loadDatafromRow(datatype, row):
    try:
        filestr = folder+'result_'+datatype+'.csv'
        f = open(filestr, "r")
        data = np.loadtxt(f, delimiter=",", skiprows=1, usecols=(row))
        return data
    except FileNotFoundError:
        print(filestr)
        print("File not found.")
        return []


# needs to be at the very end of the file
if __name__ == '__main__':
    main()
