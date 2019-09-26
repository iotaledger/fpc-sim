import numpy as np
import matplotlib.pyplot as plt

# - - - - most frequently changed output parameters - - -
# 0:a , 1:b , 2:beta , 3:k , 4:m , 5:l , 6:N , 7:p_0 , 8:q , 9:Adv_strategy
# 10:rateRandomness , 11:deltaWS , 12:gammaWS , 13:maxTermRound, 14:Type , 15:X , 16:Y
xlabel = "Proportion of adversary nodes"
xcol = 8
ylabel = "p0"
ycol = 7

# - - - - parameters unlikely to be changed - - -
zcol = 16  # last column in the csv file
folder = "data/"  # dont change this


def main():
    printData("Agreement", "Blues")
    printData("Integrity", "Oranges")
    printData("Termination", "Greens")


def printData(type, cmap):

    x = loadDatafromRow(type+"Rate", xcol)
    y = loadDatafromRow(type+"Rate", ycol)
    z = loadDatafromRow(type+"Rate", zcol)
    X, Y, Z = get2DMatrices(x, y, z)
    # print("X=", X)
    # print("Y=", Y)
    # print(type)
    Z[Z < 0.8] = 0.8
    # print("Z=", Z)

    # the option figsize  can take very long
    fig = plt.figure()
    ax = plt.axes()
    # cax = plt.contourf(X, Y, Z,  levels=np.linspace(0.8, 1, num=21), cmap=cmap)
    cax = ax.imshow(np.rot90(Z), extent=[x[0], x[-1], y[0], y[-1]],
                    cmap=cmap, vmin=0.8, vmax=1, aspect='auto')
    # cax = ax.imshow(np.rot90(Z), extent=[x[0], x[-1], .5, 1],
    #                 cmap=cmap, vmin=0.8, vmax=1, aspect='auto')
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    cbar = fig.colorbar(cax, ticks=[0.8, 0.9, 1])
    cbar.ax.set_yticklabels(['<0.8', '0.9', '1.0'])
    plt.title(type+"Rate")
    plt.savefig(folder+"plot2d_"+type+"Rate.eps", format='eps')
    plt.clf()


def get2DMatrices(x, y, z):
    xunique = np.unique(np.sort(x))
    yunique = np.unique(np.sort(y))
    xlen = len(xunique)
    ylen = len(yunique)
    X = np.empty((xlen, ylen))
    Y = np.empty((xlen, ylen))
    Z = np.empty((xlen, ylen))

    ix = getIndexList(x, xunique)
    iy = getIndexList(y, yunique)

    for i in range(len(x)):
        i1 = int(ix[i])
        i2 = int(iy[i])
        X[i1, i2] = x[i]
        Y[i1, i2] = y[i]
        Z[i1, i2] = z[i]
    return X, Y, Z


def getIndexList(x, uniquex):
    y = np.empty(len(x))
    for i in range(len(x)):
        y[i] = int(np.where(uniquex == x[i])[0][0])
    return y


def sort3vecs(x, y, z):
    i = np.argsort(x)
    x = x[i]
    y = y[i]
    z = z[i]
    return x, y, z


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
