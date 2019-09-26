import numpy as np
import matplotlib.pyplot as plt

# - - - - most frequently changed output parameters - - -
# 0:a , 1:b , 2:beta , 3:k , 4:m , 5:l , 6:N , 7:p_0 , 8:q , 9:Adv_strategy
# 10:rateRandomness , 11:deltaWS , 12:gammaWS , 13:maxTermRound, 14:Type , 15:X , 16:Y
xcol = 11
xlabel = "Proportion of queryable network $\\delta$"
xlim = [.0, 1.]
#xlim = [1, 100]
# xscale = 'log'
xscale = 'linear'

# - - - - parameters unlikely to be changed - - -
ycol = 16  # last column in the csv file
folder = "data/"


def main():
    printDataRates()
    printMeanTermRound()


def printDataRates():
    fig = plt.figure()
    filename = folder+'plot_AR-IR-TR'
    partPlot("Agreement", "AgreementRate",  ".", 10, filename, "blue")
    partPlot("Integrity", "IntegrityRate",   ".", 10, filename, "orange")
    partPlot("Termination", "TerminationRate",   ".", 10, filename, "green")
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
    # if style == '':
    #     plt.plot(x, y, label=type)
    # else:
    #     plt.plot(x, y, label=type, marker=style, linestyle='none')
    plt.plot(x, y, linestyle='dashed', color=color, linewidth=1)
    plt.plot(x, y, label=type, marker=marker,
             linestyle='none', color=color, markersize=markersize)
    np.savez(filename+"_"+type, x=x, y=y)


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
